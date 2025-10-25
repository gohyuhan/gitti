package git

import (
	"fmt"
	"os/exec"
	"strings"
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
	RepoPath      string
	FilesStatus   []FileStatus
	FilesPosition map[string]int
	UpdateChannel chan string
	ErrorLog      []error
}

func InitGitFile(repoPath string, updateChannel chan string) {
	gitFiles := GitFiles{
		RepoPath:      repoPath,
		FilesStatus:   make([]FileStatus, 0),
		UpdateChannel: updateChannel,
	}
	GITFILES = &gitFiles
}

func (gf *GitFiles) GetGitFilesStatus() {
	gitArgs := []string{"status", "--porcelain"}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gf.RepoPath
	gitOutput, err := cmd.Output()
	if err != nil {
		gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("[GIT FILES ERROR]: %w", err))
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
		fileName := strings.TrimSpace(file[3:])

		// check if this was also in the previsou list before any update to the list and retrieve back the SelectedForStage info
		fileIndex, exist := gf.FilesPosition[fileName]
		selectedForStage := true
		if exist {
			selectedForStage = gf.FilesStatus[fileIndex].SelectedForStage
		}

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FileName:         fileName,
			IndexState:       indexState,
			WorkTree:         worktree,
			SelectedForStage: selectedForStage,
		})
		modifiedFilesPositionHashmap[fileName] = index
	}
	gf.FilesPosition = modifiedFilesPositionHashmap
	gf.FilesStatus = modifiedFilesStatus
}

// get the file diff content
func (gf *GitFiles) GetFilesDiffInfo(fileName string) []FileDiffLine {
	gitArgs := []string{"diff", "HEAD", "-U99999", fileName}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = gf.RepoPath
	gitOutput, err := cmd.Output()
	if err != nil {
		gf.ErrorLog = append(gf.ErrorLog, fmt.Errorf("[GIT FILES DIFF ERROR]: %w", err))
		return nil
	}

	fileDiffOneLineString := strings.Split(string(gitOutput), "@@")
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
	fileIndex, exist := gf.FilesPosition[fileName]
	if exist {
		if gf.FilesStatus[fileIndex].SelectedForStage {
			gf.FilesStatus[fileIndex].SelectedForStage = false
		} else {
			gf.FilesStatus[fileIndex].SelectedForStage = true
		}
		gf.UpdateChannel <- GIT_FILES_STATUS_UPDATE
	}
}
