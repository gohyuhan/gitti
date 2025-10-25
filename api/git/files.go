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
	RepoPath            string
	CurrentSelectedFile string
	FilesStatus         []FileStatus
	ErrorLog            []error
}

func InitGitFile(repoPath string) {
	gitFiles := GitFiles{
		RepoPath:            repoPath,
		FilesStatus:         make([]FileStatus, 0),
		CurrentSelectedFile: "",
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

	files := strings.SplitSeq(strings.TrimRight(string(gitOutput), "\n"), "\n")

	currentModifiedFiles := []string{}
	modifiedFilesStatus := []FileStatus{}

	for file := range files {
		if len(file) < 3 {
			continue
		}

		indexState := string(file[0])
		worktree := string(file[1])
		fileName := strings.TrimSpace(file[3:])

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FileName:         fileName,
			IndexState:       indexState,
			WorkTree:         worktree,
			SelectedForStage: true,
		})
		currentModifiedFiles = append(currentModifiedFiles, fileName)
	}

	// to reassign the current file if the new files doesn't contain the current selected GetGitFilesStatus
	if len(currentModifiedFiles) < 1 {
		gf.CurrentSelectedFile = ""
	} else {
		if !Contains(currentModifiedFiles, gf.CurrentSelectedFile) || gf.CurrentSelectedFile == "" {
			gf.CurrentSelectedFile = currentModifiedFiles[0]
		}
	}
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
