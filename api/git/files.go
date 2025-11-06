package git

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"gitti/cmd"
)

var GITFILES *GitFiles

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
	FilesStatus                 []FileStatus
	FilesPosition               map[string]int
	UpdateChannel               chan string
	FilesSelectedForStageStatus map[string]bool // centralized recording the user selection if they want to stage a file
	ErrorLog                    []error
	GitFilesMutex               sync.Mutex
}

func InitGitFile(updateChannel chan string) {
	gitFiles := GitFiles{
		FilesStatus:                 make([]FileStatus, 0),
		UpdateChannel:               updateChannel,
		FilesSelectedForStageStatus: make(map[string]bool),
	}
	GITFILES = &gitFiles
}

func (gf *GitFiles) GetGitFilesStatus() {
	gitArgs := []string{"status", "--porcelain", "-uall"}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gitOutput, err := cmd.Output()
	if err != nil {
		gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("[GIT FILES ERROR]: %w", err))
	}

	files := strings.Split(strings.TrimRight(string(gitOutput), "\n"), "\n")

	modifiedFilesStatus := []FileStatus{}
	modifiedFilesPositionHashmap := make(map[string]int)

	gf.GitFilesMutex.Lock()
	for index, file := range files {
		if len(file) < 3 {
			continue
		}

		indexState := string(file[0])
		worktree := string(file[1])
		fileName := strings.TrimSpace(file[3:])

		// check if this was also in the previsou list before any update to the list and retrieve back the SelectedForStage info
		_, exist := gf.FilesSelectedForStageStatus[fileName]
		if !exist {
			gf.FilesSelectedForStageStatus[fileName] = true
		}

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FileName:         fileName,
			IndexState:       indexState,
			WorkTree:         worktree,
			SelectedForStage: gf.FilesSelectedForStageStatus[fileName],
		})
		modifiedFilesPositionHashmap[fileName] = index
	}

	gf.FilesPosition = modifiedFilesPositionHashmap
	gf.FilesStatus = modifiedFilesStatus
	gf.GitFilesMutex.Unlock()
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
				gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("[GIT FILES DIFF ERROR]: %w", err))
				return nil
			}
		} else {
			gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("[GIT FILES DIFF ERROR]: %w", err))
			return nil
		}
	}

	fileDiffOneLineString := strings.SplitN(string(gitOutput), "@@", 3)
	if len(fileDiffOneLineString) < 3 {
		gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("There is no diff for the selected file or the file format is not supported for preview"))
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

func (gf *GitFiles) ToggleFilesStageStatus(fileName string) {
	gf.GitFilesMutex.Lock()
	_, fileStatusExist := gf.FilesSelectedForStageStatus[fileName]
	fileIndex, fileIndexExist := gf.FilesPosition[fileName]
	if fileIndexExist && fileStatusExist {
		gf.FilesSelectedForStageStatus[fileName] = !gf.FilesSelectedForStageStatus[fileName]
		gf.FilesStatus[fileIndex].SelectedForStage = gf.FilesSelectedForStageStatus[fileName]
		gf.UpdateChannel <- GIT_FILES_STATUS_UPDATE
	}
	gf.GitFilesMutex.Unlock()
}
