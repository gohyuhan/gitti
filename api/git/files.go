package git

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

// const (
// 	AddLine        = "ADDLINE"
// 	RemoveLine     = "REMOVELINE"
// 	UnModifiedLine = "UNMODIFIEDLINE"
// )

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
func (gf *GitFiles) GetFilesDiffInfo(ctx context.Context, fileStatus FileStatus) []string {
	filePathName := fileStatus.FilePathname
	if fileStatus.IndexState == "R" || fileStatus.IndexState == "C" {
		if strings.Contains(filePathName, "->") {
			parts := strings.Split(filePathName, "->")
			if len(parts) >= 2 {
				filePathName = strings.TrimSpace(parts[1])
			}
		}
	}
	gitArgs := []string{"diff", "HEAD", "--", filePathName}
	// the file is untracked
	isNewFile := fileStatus.WorkTree == "?" ||
		fileStatus.IndexState == "?" ||
		fileStatus.IndexState == "A" ||
		(fileStatus.IndexState == "U" && fileStatus.WorkTree == "A")

	if isNewFile {
		// empty file for git diff --no-index to compares two arbitrary files outside the Git index.
		nullFile := "/dev/null"
		if runtime.GOOS == "windows" {
			nullFile = "NUL"
		}
		gitArgs = []string{"diff", "--no-index", nullFile, "--", filePathName}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	gitOutput, err := cmdExecutor.Output()
	if err != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gf.errorLog = append(gf.errorLog, fmt.Errorf("[FILE DIFF OPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err()))
			return nil
		}
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

	fileDiffLines := processGeneralGitOpsOutputIntoStringArray(gitOutput)
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
		// "old -> new" format for both Renamed (R) and Copied (C)
		// This covers IndexState R/C and the rare WorkTree R/C
		if strings.Contains(filePathName, "->") &&
			(file.IndexState == "R" || file.IndexState == "C" || file.WorkTree == "R" || file.WorkTree == "C") {
			filePathName = strings.TrimSpace(strings.Split(filePathName, "->")[1])
		}

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
	needToRunExecutor := true
	var gitArgs []string

	fileIndex, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		file := gf.filesStatus[fileIndex]
		filePathName = file.FilePathname
		// Store the full name (e.g. "old -> new") for the rename logic later
		fullFilePathName := file.FilePathname

		// "old -> new" format for both Renamed (R) and Copied (C)
		// This covers IndexState R/C and the rare WorkTree R/C
		if strings.Contains(filePathName, "->") &&
			(file.IndexState == "R" || file.IndexState == "C" || file.WorkTree == "R" || file.WorkTree == "C") {
			filePathName = strings.TrimSpace(strings.Split(filePathName, "->")[1])
		}

		switch discardType {
		case DISCARDWHOLE:
			gitArgs = []string{"checkout", "HEAD", "--", filePathName}
		case DISCARDUNSTAGE:
			gitArgs = []string{"checkout", "--", filePathName}
		case DISCARDUNTRACKED:
			gitArgs = []string{"clean", "-f", "--", filePathName}
			// although they are in worktree, they are actually tracked, therefore we need to use git rm -f <filename>
			if file.WorkTree == "A" || file.WorkTree == "C" || file.WorkTree == "R" {
				gitArgs = []string{"rm", "-f", filePathName}
			}
			// we are refetching it actively here is because the clean doesn't trigger any write in .git folder
			// and therefore will not trigger the watcher event driven fetch for file status, so we trigger a fetch here
			// to prevent a "lag" in the UI
			needFilesStatusRefetch = true
		case DISCARDNEWLYADDEDORCOPIED:
			gitArgs = []string{"rm", "-f", filePathName}
		case DISCARDANDREVERTRENAME:
			needToRunExecutor = false
			oldFilePathName := strings.TrimSpace(strings.Split(fullFilePathName, "->")[0])
			newFilePathName := strings.TrimSpace(strings.Split(fullFilePathName, "->")[1])

			// retrieve back the original file
			gitArgs = []string{"reset", "--", oldFilePathName}
			oldFileResetCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			oldFileResetCmdExecutor.Run()

			gitArgs = []string{"checkout", "--", oldFilePathName}
			oldFileRevertCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			oldFileRevertCmdExecutor.Run()

			// revert and remove the "newly named" file
			gitArgs = []string{"reset", "--", newFilePathName}
			newFileResetCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			newFileResetCmdExecutor.Run()

			gitArgs = []string{"clean", "-f", "--", newFilePathName}
			newFileDiscardCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			newFileDiscardCmdExecutor.Run()

			needFilesStatusRefetch = true
		}

		if needToRunExecutor {
			changesDiscardCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
			changesDiscardCmdExecutor.Run()
		}

		if needFilesStatusRefetch {
			go func() {
				gf.GetGitFilesStatus()
				gf.updateChannel <- GIT_FILES_STATUS_UPDATE
			}()
		}
	}
}
