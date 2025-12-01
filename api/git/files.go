package git

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

const (
	AddLine        = "ADDLINE"
	RemoveLine     = "REMOVELINE"
	UnModifiedLine = "UNMODIFIEDLINE"
)

const (
	DISCARDWHOLE      = "DISCARDWHOLE"
	DISCARDUNSTAGE    = "DISCARDUNSTAGE"
	DISCARDUNTRACKED  = "DISCARDUNTRACKED"
	DISCARDNEWLYADDED = "DISCARDNEWLYADDED"
)

type FileStatus struct {
	FilePathname string
	IndexState   string
	WorkTree     string
}

type GitFiles struct {
	filesStatus    []FileStatus
	filesPosition  map[string]int
	errorLog       []error
	gitProcessLock *GitProcessLock
	updateChannel  chan string
}

func InitGitFile(updateChannel chan string, gitProcessLock *GitProcessLock) *GitFiles {
	gitFiles := GitFiles{
		filesStatus:    make([]FileStatus, 0),
		gitProcessLock: gitProcessLock,
		updateChannel:  updateChannel,
	}
	return &gitFiles
}

// ----------------------------------
//
//	Return filesStatus
//
// ----------------------------------
func (gf *GitFiles) FilesStatus() []FileStatus {
	copied := make([]FileStatus, len(gf.filesStatus))
	copy(copied, gf.filesStatus)
	return copied
}

// ----------------------------------
//
//		Retrieve File Status
//	 * Passive, this should only be trigger by system
//
// ----------------------------------
func (gf *GitFiles) GetGitFilesStatus() {
	gitArgs := []string{"status", "--porcelain", "--untracked-files=all"}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gitOutput, err := cmdExecutor.Output()
	if err != nil {
		gf.errorLog = append(gf.errorLog, fmt.Errorf("[GIT FILES ERROR]: %w", err))
	}

	files := strings.Split(strings.TrimRight(string(gitOutput), "\n"), "\n")

	modifiedFilesStatus := []FileStatus{}
	modifiedFilesPositionHashmap := make(map[string]int)

	for index, file := range files {
		if len(file) < 3 {
			continue
		}

		indexState := string(file[0])
		worktree := string(file[1])
		filePathName := strings.TrimSpace(file[3:])

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FilePathname: filePathName,
			IndexState:   indexState,
			WorkTree:     worktree,
		})
		modifiedFilesPositionHashmap[filePathName] = index
	}

	gf.filesPosition = modifiedFilesPositionHashmap
	gf.filesStatus = modifiedFilesStatus
}

// get the file diff content
func (gf *GitFiles) GetFilesDiffInfo(fileStatus FileStatus) []string {
	filePathName := fileStatus.FilePathname
	if fileStatus.IndexState == "R" {
		filePathName = strings.TrimSpace(strings.Split(filePathName, "->")[1])
	}
	gitArgs := []string{"diff", "HEAD", "--", filePathName}
	// the file is untracked
	if fileStatus.WorkTree == "?" || fileStatus.IndexState == "?" || fileStatus.IndexState == "A" {
		// empty file for git diff --no-index to compares two arbitrary files outside the Git index.
		nullFile := "/dev/null"
		if runtime.GOOS == "windows" {
			nullFile = "NUL"
		}
		gitArgs = []string{"diff", "--no-index", nullFile, "--", filePathName}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, true)
	gitOutput, err := cmdExecutor.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			if exitError.ExitCode() != 1 {
				gf.errorLog = append(gf.errorLog, fmt.Errorf("[GIT FILES DIFF ERROR]: %w", err))
				return nil
			}
		} else {
			gf.errorLog = append(gf.errorLog, fmt.Errorf("[GIT FILES DIFF ERROR]: %w", err))
			return nil
		}
	}

	fileDiffLines := strings.Split(strings.TrimSpace(string(gitOutput)), "\n")
	return fileDiffLines
}

func (gf *GitFiles) StageOrUnstageFile(filePathName string) {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	fileIndex, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		file := gf.filesStatus[fileIndex]
		var gitArgs []string
		if file.IndexState == "?" && file.WorkTree == "?" {
			// not tracked
			gitArgs = []string{"add", "--", filePathName}
			stageCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			stageCmdExecutor.Run()
		} else if file.IndexState != " " && file.WorkTree != " " {
			// staged but have modification later
			gitArgs = []string{"add", "--", filePathName}
			stageCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			stageCmdExecutor.Run()
		} else if file.IndexState != " " && file.WorkTree == " " {
			// staged and no latest modification, so we need to unstage it or revert back
			gitArgs = []string{"reset", "--", filePathName}
			if file.IndexState == "A" {
				gitArgs = []string{"rm", "--cached", "--force", "--", filePathName}
			}
			unstageCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			unstageCmd.Run()
		} else if file.IndexState == " " && file.WorkTree != " " {
			// tracked but not staged
			gitArgs = []string{"add", "--", filePathName}
			stageCmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			stageCmd.Run()
		}
	}
}

func (gf *GitFiles) StageAllChanges() {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"add", "."}
	stageCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stageCmdExecutor.Run()
}

func (gf *GitFiles) UnstageAllChanges() {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	gitArgs := []string{"reset"}
	stageCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	stageCmdExecutor.Run()

}

// ----------------------------------
//
//	Discard File changes
//
// ----------------------------------
func (gf *GitFiles) DiscardFileChanges(filePathName string, discardType string) {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	needFilesStatusRefetch := false
	var gitArgs []string
	switch discardType {
	case DISCARDWHOLE:
		gitArgs = []string{"checkout", "HEAD", "--", filePathName}
	case DISCARDUNSTAGE:
		gitArgs = []string{"checkout", "--", filePathName}
	case DISCARDUNTRACKED:
		gitArgs = []string{"clean", "-f", "--", filePathName}
		// we are refetching it actively here is because the clean doesn't trigger any write in .git folder
		// and therefore will not trigger the watcher event driven fetch for file status, so we trigger a fetch here
		// to prevent a "lag" in the UI
		needFilesStatusRefetch = true
	case DISCARDNEWLYADDED:
		gitArgs = []string{"rm", "-f", filePathName}
	}
	changesDiscardCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	changesDiscardCmdExecutor.Run()

	if needFilesStatusRefetch {
		go func() {
			gf.GetGitFilesStatus()
			gf.updateChannel <- GIT_FILES_STATUS_UPDATE
		}()
	}
}
