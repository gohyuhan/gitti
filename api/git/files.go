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
	FileName         string
	IndexState       string
	WorkTree         string
	SelectedForStage bool
}

type FileDiffLine struct {
	Line string
	Type string
}

type GitFiles struct {
	filesStatus                 []FileStatus
	filesPosition               map[string]int
	updateChannel               chan string
	filesSelectedForStageStatus map[string]bool // centralized recording the user selection if they want to stage a file
	errorLog                    []error
	gitFilesMutex               sync.Mutex
}

func InitGitFile(updateChannel chan string) *GitFiles {
	gitFiles := GitFiles{
		filesStatus:                 make([]FileStatus, 0),
		updateChannel:               updateChannel,
		filesSelectedForStageStatus: make(map[string]bool),
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

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
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
		fileName := strings.TrimSpace(file[3:])

		// check if this was also in the previous list before any update to the list and retrieve back the SelectedForStage info
		_, exist := gf.filesSelectedForStageStatus[fileName]
		if !exist {
			gf.filesSelectedForStageStatus[fileName] = true
		}

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FileName:         fileName,
			IndexState:       indexState,
			WorkTree:         worktree,
			SelectedForStage: gf.filesSelectedForStageStatus[fileName],
		})
		modifiedFilesPositionHashmap[fileName] = index
	}

	gf.filesPosition = modifiedFilesPositionHashmap
	gf.filesStatus = modifiedFilesStatus
	gf.gitFilesMutex.Unlock()
}

// get the file diff content
func (gf *GitFiles) GetFilesDiffInfo(fileStatus FileStatus) []FileDiffLine {
	gitArgs := []string{"diff", "HEAD", "--diff-filter=ADM", "-U99999", "--", fileStatus.FileName}
	// the file is untracked
	if fileStatus.WorkTree == "?" || fileStatus.IndexState == "?" {
		// empty file for git diff --no-index to compares two arbitrary files outside the Git index.
		nullFile := "/dev/null"
		if runtime.GOOS == "windows" {
			nullFile = "NUL"
		}
		gitArgs = []string{"diff", "--no-index", "-U99999", nullFile, "--", fileStatus.FileName}
	}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
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

func (gf *GitFiles) GetSelectedForStageFiles() []string {
	var stagedFiles []string
	gf.gitFilesMutex.Lock()
	for _, file := range gf.filesStatus {
		if file.SelectedForStage {
			stagedFiles = append(stagedFiles, file.FileName)
		}
	}
	gf.gitFilesMutex.Unlock()
	return stagedFiles
}

func (gf *GitFiles) ToggleFilesStageStatus(fileName string) {
	gf.gitFilesMutex.Lock()
	_, fileStatusExist := gf.filesSelectedForStageStatus[fileName]
	fileIndex, fileIndexExist := gf.filesPosition[fileName]
	if fileIndexExist && fileStatusExist {
		gf.filesSelectedForStageStatus[fileName] = !gf.filesSelectedForStageStatus[fileName]
		gf.filesStatus[fileIndex].SelectedForStage = gf.filesSelectedForStageStatus[fileName]
		gf.updateChannel <- GIT_FILES_STATUS_UPDATE
	}
	gf.gitFilesMutex.Unlock()
}
