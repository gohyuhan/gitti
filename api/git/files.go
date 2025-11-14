package git

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"gitti/cmd"
)

const (
	AddLine        = "ADDLINE"
	RemoveLine     = "REMOVELINE"
	UnModifiedLine = "UNMODIFIEDLINE"
)

type FileStatus struct {
	FilePathname string
	IndexState   string
	WorkTree     string
}

type FileDiffLine struct {
	Line string
	Type string
}

type GitFiles struct {
	filesStatus    []FileStatus
	filesPosition  map[string]int
	errorLog       []error
	gitFilesMutex  sync.Mutex
	gitProcessLock *GitProcessLock
}

func InitGitFile(gitProcessLock *GitProcessLock) *GitFiles {
	gitFiles := GitFiles{
		filesStatus:    make([]FileStatus, 0),
		gitProcessLock: gitProcessLock,
	}
	return &gitFiles
}

// ----------------------------------
//
//	Return filesStatus
//
// ----------------------------------
func (gf *GitFiles) FilesStatus() []FileStatus {
	gf.gitFilesMutex.Lock()
	defer gf.gitFilesMutex.Unlock()

	copied := make([]FileStatus, len(gf.filesStatus))
	copy(copied, gf.filesStatus)
	return copied
}

// ----------------------------------
//
//		Retrieve File Status
//	 * Passive, this should onyl be trigger by system
//
// ----------------------------------
func (gf *GitFiles) GetGitFilesStatus() {
	gitArgs := []string{"status", "--porcelain", "-uall"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	gitOutput, err := cmd.Output()
	if err != nil {
		gf.errorLog = append(gf.errorLog, fmt.Errorf("[GIT FILES ERROR]: %w", err))
	}

	files := strings.Split(strings.TrimRight(string(gitOutput), "\n"), "\n")

	modifiedFilesStatus := []FileStatus{}
	modifiedFilesPositionHashmap := make(map[string]int)

	gf.gitFilesMutex.Lock()
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
	gf.gitFilesMutex.Unlock()
}

// get the file diff content
func (gf *GitFiles) GetFilesDiffInfo(fileStatus FileStatus) []FileDiffLine {
	gitArgs := []string{"diff", "HEAD", "--diff-filter=ADM", "-U99999", "--", fileStatus.FilePathname}
	// the file is untracked
	if fileStatus.WorkTree == "?" || fileStatus.IndexState == "?" {
		// empty file for git diff --no-index to compares two arbitrary files outside the Git index.
		nullFile := "/dev/null"
		if runtime.GOOS == "windows" {
			nullFile = "NUL"
		}
		gitArgs = []string{"diff", "--no-index", "-U99999", nullFile, "--", fileStatus.FilePathname}
	}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	gitOutput, err := cmd.Output()
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

	fileDiffOneLineString := strings.SplitN(string(gitOutput), "@@", 3)
	if len(fileDiffOneLineString) < 3 {
		gf.errorLog = append(gf.errorLog, fmt.Errorf("There is no diff for the selected file or the file format is not supported for preview"))
		return nil
	}
	fileDiffLines := strings.SplitSeq(strings.TrimSpace(fileDiffOneLineString[2]), "\n")
	fileDiff := []FileDiffLine{}

	for Line := range fileDiffLines {
		fileDiffLine := FileDiffLine{
			Line: Line,
			Type: UnModifiedLine,
		}
		if strings.HasPrefix(Line, "-") {
			fileDiffLine.Type = RemoveLine
		} else if strings.HasPrefix(Line, "+") {
			fileDiffLine.Type = AddLine
		}

		fileDiff = append(fileDiff, fileDiffLine)
	}
	return fileDiff
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
			stageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
			stageCmd.Run()
		} else if file.IndexState != " " && file.WorkTree != " " {
			// staged but have modification later
			gitArgs = []string{"add", "--", filePathName}
			stageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
			stageCmd.Run()
		} else if file.IndexState != " " && file.WorkTree == " " {
			// staged and no latest modification, so we need to unstage it or revert back
			gitArgs = []string{"reset", "HEAD", "--", filePathName}
			if file.IndexState == "A" {
				gitArgs = []string{"rm", "--cached", "--force", "--", filePathName}
			}
			unstageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
			unstageCmd.Run()
		} else if file.IndexState == " " && file.WorkTree != " " {
			// tracked but not staged
			gitArgs = []string{"add", "--", filePathName}
			stageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
			stageCmd.Run()
		}
	}
}

func (gf *GitFiles) StageAllChanges() {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	var gitArgs []string

	gitArgs = []string{"add", "."}
	stageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	stageCmd.Run()
}

func (gf *GitFiles) UnstageAllChanges() {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	var gitArgs []string

	gitArgs = []string{"reset", "HEAD"}
	stageCmd := cmd.GittiCmd.RunGitCmd(gitArgs, false)
	stageCmd.Run()

}
