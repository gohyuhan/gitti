package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/gohyuhan/gitti/executor"
)

type FileStatus struct {
	FilePathname string
	IndexState   string
	WorkTree     string
	HasConflict  bool
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
//	Retrieve File Status
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
		hasConflict := isFilesInConflictState(indexState, worktree)

		modifiedFilesStatus = append(modifiedFilesStatus, FileStatus{
			FilePathname: filePathName,
			IndexState:   indexState,
			WorkTree:     worktree,
			HasConflict:  hasConflict,
		})
		modifiedFilesPositionHashmap[filePathName] = index
	}

	gf.filesPosition = modifiedFilesPositionHashmap
	gf.filesStatus = modifiedFilesStatus
}

// get the file diff content
func (gf *GitFiles) GetFilesDiffInfo(ctx context.Context, fileStatus FileStatus, DiffType string) []string {
	filePathName := fileStatus.FilePathname
	if fileStatus.IndexState == "R" || fileStatus.IndexState == "C" {
		if strings.Contains(filePathName, "->") {
			parts := strings.Split(filePathName, "->")
			if len(parts) >= 2 {
				filePathName = strings.TrimSpace(parts[1])
			}
		}
	}
	var gitArgs []string
	switch DiffType {
	case GETSTAGEDDIFF:
		gitArgs = []string{"diff", "--cached", "--", filePathName}
	case GETUNSTAGEDDIFF:
		gitArgs = []string{"diff", "--", filePathName}
	case GETCOMBINEDDIFF:
		gitArgs = []string{"diff", "HEAD", "--", filePathName}
	}
	// the file is untracked
	isNewFile := fileStatus.WorkTree == "?" ||
		fileStatus.IndexState == "?" ||
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
//	Stage line
//
// ----------------------------------
func (gf *GitFiles) StageLine(filePathName string, diffContentStringArray []string, startFromIndex int, stageLineIndex int) {
	// Acquire Git process lock to ensure no other Git operations are running concurrently.
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	// Calculate the index of the line relative to the start of the diff chunk
	actualStageLineIndex := stageLineIndex - startFromIndex

	// Bounds check to prevent panic and ensure we only process modified lines (+/-)
	if actualStageLineIndex < 0 || actualStageLineIndex >= len(diffContentStringArray) {
		return
	}
	selectedLine := diffContentStringArray[actualStageLineIndex]
	isAddition := strings.HasPrefix(selectedLine, "+") && !strings.HasPrefix(selectedLine, "+++")
	isDeletion := strings.HasPrefix(selectedLine, "-") && !strings.HasPrefix(selectedLine, "---")
	if !isAddition && !isDeletion {
		return
	}

	fileIndex, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		file := gf.filesStatus[fileIndex]
		// "old -> new" format for both Renamed (R) and Copied (C)
		// This covers IndexState R/C and the rare WorkTree R/C
		if strings.Contains(filePathName, "->") &&
			(file.IndexState == "R" || file.IndexState == "C" || file.WorkTree == "R" || file.WorkTree == "C") {
			filePathName = strings.TrimSpace(strings.Split(filePathName, "->")[1])
		}

		// Create a temporary patch file.
		// This file will hold the unified diff content for the specific line we want to stage.
		tempPatchFile, err := os.CreateTemp("", "gitti-patch-stage-*")
		if err != nil {
			return
		}
		// Ensure the file is cleaned up (deleted) after function exits
		defer os.Remove(tempPatchFile.Name())

		var gitArgs []string
		// Generate and write the patch content.
		// We generate a custom patch that isolates the specific line change and write it to the temp file.
		_, writeErr := tempPatchFile.WriteString(generateStageLinePatchString(diffContentStringArray, actualStageLineIndex, filePathName))
		if writeErr != nil {
			tempPatchFile.Close() // Close before return
			return
		}
		// Close the file descriptor so Git can read it safely
		closeErr := tempPatchFile.Close()
		if closeErr != nil {
			return
		}

		// Apply the patch to the index.
		// We use `git apply --cached` to apply the patch to the staging area only.
		// --recount: Allow overlapping hunks if offsets change.
		// --whitespace=nowarn: Ignore whitespace warnings.
		gitArgs = []string{"apply", "--cached", "--recount", "--whitespace=nowarn", tempPatchFile.Name()}
		stageLineCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
		stageLineCmdExecutor.Run()

		// Update file status to reflect changes
		gf.GetGitFilesStatus()
		gf.updateChannel <- GIT_EDIT_LINE_DETAILS_AND_FILES_UPDATE
	}
}

// generateStageLinePatchString creates a patch string for a single line from a larger diff.
// It keeps the target line as a change, converts other changes in the same context to "context lines",
// and preserves existing context lines. This effectively isolates the chosen line for staging/unstaging.
func generateStageLinePatchString(diffContentStringArray []string, actualStageLineIndex int, filePathName string) string {
	var stageLinePatchString strings.Builder
	// We need to track if the previous line was skipped to handle "\ No newline" markers correctly.
	lastLineWasSkipped := false

	for index, diffLine := range diffContentStringArray {
		if utf8.RuneCountInString(diffLine) > 0 {
			// ---------------------------------------------------------
			// HANDLE SPECIAL MARKERS (\ No newline)
			// ---------------------------------------------------------
			// If the previous line was skipped (e.g. an unstaged addition),
			// we likely need to skip this marker too, otherwise it might attach
			// to a context line that shouldn't have it.
			if strings.HasPrefix(diffLine, "\\ No newline") {
				if lastLineWasSkipped {
					continue
				}
				stageLinePatchString.WriteString(diffLine)
				stageLinePatchString.WriteString("\n")
				continue
			}
			if strings.HasPrefix(diffLine, "deleted file mode") {
				continue
			}
			if strings.HasPrefix(diffLine, "+++ /dev/null") {
				stageLinePatchString.WriteString("+++ b/")
				stageLinePatchString.WriteString(filePathName)
				stageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
				continue
			}
			// If this is the line we want to stage, include it exactly as is (e.g., "+ line" or "- line").
			if index == actualStageLineIndex {
				stageLinePatchString.WriteString(diffLine)
				stageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
				continue
			}
			// For other lines in the diff hunk:
			if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "---") {
				// If it's a deletion line ("- content"), convert it to a context line ("  content").
				// This means "treat this as existing unchanged content" in the patch context.
				if utf8.RuneCountInString(diffLine) > 1 {
					stageLinePatchString.WriteString(" ")
					stageLinePatchString.WriteString(diffLine[1:])
					stageLinePatchString.WriteString("\n")
				} else {
					stageLinePatchString.WriteString(" \n")
				}
				lastLineWasSkipped = false
			} else if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++") {
				// If it's an addition line ("+ content"), skip it.
				// We don't want to include other additions in this specific single-line patch.
				lastLineWasSkipped = true
				continue
			} else {
				// For context lines (starting with space) or header lines, keep them as is.
				stageLinePatchString.WriteString(diffLine)
				stageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
			}
		} else {
			stageLinePatchString.WriteString("\n")
			lastLineWasSkipped = false
		}
	}
	return stageLinePatchString.String()
}

// ----------------------------------
//
//	Unstage line
//
// ----------------------------------
func (gf *GitFiles) UnstageLine(filePathName string, diffContentStringArray []string, startFromIndex int, unStageLineIndex int) {
	// Acquire Git process lock.
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	// Calculate the index of the line relative to the start of the diff chunk
	actualUnStageLineIndex := unStageLineIndex - startFromIndex

	// Bounds check to prevent panic and ensure we only process modified lines (+/-)
	if actualUnStageLineIndex < 0 || actualUnStageLineIndex >= len(diffContentStringArray) {
		return
	}
	selectedLine := diffContentStringArray[actualUnStageLineIndex]
	isAddition := strings.HasPrefix(selectedLine, "+") && !strings.HasPrefix(selectedLine, "+++")
	isDeletion := strings.HasPrefix(selectedLine, "-") && !strings.HasPrefix(selectedLine, "---")
	if !isAddition && !isDeletion {
		return
	}

	fileIndex, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		file := gf.filesStatus[fileIndex]
		// "old -> new" format for both Renamed (R) and Copied (C)
		// This covers IndexState R/C and the rare WorkTree R/C
		if strings.Contains(filePathName, "->") &&
			(file.IndexState == "R" || file.IndexState == "C" || file.WorkTree == "R" || file.WorkTree == "C") {
			filePathName = strings.TrimSpace(strings.Split(filePathName, "->")[1])
		}

		if file.IndexState == "?" && file.WorkTree == "?" {
			// not tracked
			return
		}

		// Create a temporary patch file.
		tempPatchFile, err := os.CreateTemp("", "gitti-patch-unstage-*")
		if err != nil {
			return
		}
		// Ensure the file is cleaned up (deleted) after function exits
		defer os.Remove(tempPatchFile.Name())

		var gitArgs []string
		// Generate and write the patch content.
		// We reuse generateStageLinePatchString because the patch structure is the same.
		// The difference is in how we apply it (using --reverse).
		_, writeErr := tempPatchFile.WriteString(generateUnstageLinePatchString(diffContentStringArray, actualUnStageLineIndex, filePathName))
		if writeErr != nil {
			tempPatchFile.Close() // Close before return
			return
		}
		// Close the file descriptor so Git can read it safely
		closeErr := tempPatchFile.Close()
		if closeErr != nil {
			return
		}

		// Apply the patch in reverse to the index.
		// `git apply --reverse` effectively undoes the change described in the patch.
		// Since the patch describes adding/removing the line, reversing it unstages that change.
		gitArgs = []string{"apply", "--cached", "--recount", "--reverse", "--whitespace=nowarn", tempPatchFile.Name()}
		unStageLineCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
		unStageLineCmdExecutor.Run()

		// Update file status to reflect changes
		gf.GetGitFilesStatus()
		gf.updateChannel <- GIT_EDIT_LINE_DETAILS_AND_FILES_UPDATE
	}
}

func generateUnstageLinePatchString(diffContentStringArray []string, actualUnStageLineIndex int, filePathName string) string {
	var unStageLinePatchString strings.Builder
	lastLineWasSkipped := false

	for index, diffLine := range diffContentStringArray {
		if utf8.RuneCountInString(diffLine) > 0 {
			// ---------------------------------------------------------
			// HANDLE SPECIAL MARKERS (\ No newline)
			// ---------------------------------------------------------
			if strings.HasPrefix(diffLine, "\\ No newline") {
				if lastLineWasSkipped {
					continue
				}
				unStageLinePatchString.WriteString(diffLine)
				unStageLinePatchString.WriteString("\n")
				continue
			}
			// This is the specific line the user wants to UNSTAGE.
			// Keep it exactly as is (+ or -). The --reverse flag in the command will flip it.
			if index == actualUnStageLineIndex {
				unStageLinePatchString.WriteString(diffLine)
				unStageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
				continue
			}

			if strings.HasPrefix(diffLine, "new file mode") {
				continue
			}
			if strings.HasPrefix(diffLine, "--- /dev/null") {
				unStageLinePatchString.WriteString("--- a/")
				unStageLinePatchString.WriteString(filePathName)
				unStageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
				continue
			}

			// Handle unselected lines
			if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++") {
				// This is a staged addition. It IS currently in the Index.
				// We convert it to a context line (starts with space) so Git can anchor the patch.
				if utf8.RuneCountInString(diffLine) > 1 {
					unStageLinePatchString.WriteString(" ")
					unStageLinePatchString.WriteString(diffLine[1:])
					unStageLinePatchString.WriteString("\n")
				} else {
					unStageLinePatchString.WriteString(" \n")
				}
				lastLineWasSkipped = false
			} else if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "---") {
				// This is a staged deletion. It IS NOT in the Index anymore.
				// We must discard it entirely or the patch will fail to find the anchor.
				lastLineWasSkipped = true
				continue

			} else {
				// Keep Headers (diff, index, ---, +++, @@) and existing Context lines.
				unStageLinePatchString.WriteString(diffLine)
				unStageLinePatchString.WriteString("\n")
				lastLineWasSkipped = false
			}
		} else {
			// Handle empty lines in the source code
			unStageLinePatchString.WriteString("\n")
			lastLineWasSkipped = false
		}
	}
	return unStageLinePatchString.String()
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
		if file.HasConflict {
			return
		}
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

func (gf *GitFiles) GitResolveConflict(filePathName string, resolveType string) {
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	var gitArgs []string

	fileIndex, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		file := gf.filesStatus[fileIndex]
		if !file.HasConflict {
			return
		}
		filePathName = file.FilePathname
		switch resolveType {
		case RESETCONFLICT:
			gitArgs = []string{"checkout", "-m", "--", filePathName}
		case CONFLICTACCEPTOURSCHANGES:
			gitArgs = []string{"checkout", "--ours", "--", filePathName}
		case CONFLICTACCEPTTHEIRSCHANGES:
			gitArgs = []string{"checkout", "--theirs", "--", filePathName}
		}
		changesDiscardCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
		changesDiscardCmdExecutor.Run()
	}
}

// ----------------------------------
//
//		       Discard Line Change
//	  * GitDiscardFileLineChange discards a specific line change from either the Index (unstage) or the Worktree (discard).
//	    It works by generating a "reverse patch" for that specific line and applying it via 'git apply'.
//
// ----------------------------------
func (gf *GitFiles) GitDiscardFileLineChange(filePathName string, diffContentStringArray []string, startFromIndex int, discardLineIndex int, stageStatus string) {
	// Acquire Git process lock to ensure no other Git operations are running concurrently.
	if !gf.gitProcessLock.CanProceedWithGitOps() {
		return
	}
	defer gf.gitProcessLock.ReleaseGitOpsLock()

	// Calculate the index of the line relative to the start of the diff chunk
	actualDiscardLineIndex := discardLineIndex - startFromIndex

	// Bounds check to prevent panic and ensure we only process modified lines (+/-)
	if actualDiscardLineIndex < 0 || actualDiscardLineIndex >= len(diffContentStringArray) {
		return
	}
	selectedLine := diffContentStringArray[actualDiscardLineIndex]
	isAddition := strings.HasPrefix(selectedLine, "+") && !strings.HasPrefix(selectedLine, "+++")
	isDeletion := strings.HasPrefix(selectedLine, "-") && !strings.HasPrefix(selectedLine, "---")
	if !isAddition && !isDeletion {
		return
	}

	_, fileIndexExist := gf.filesPosition[filePathName]
	if fileIndexExist {
		// Create a temporary patch file.
		// This file will hold the unified diff content for the specific line we want to reverse.
		tempPatchFile, err := os.CreateTemp("", "gitti-patch-discard-line-change-*")
		if err != nil {
			return
		}
		// Ensure the file is cleaned up (deleted) after function exits
		defer os.Remove(tempPatchFile.Name())

		var gitArgs []string
		var writeErr error

		// Generate and write the patch content.
		// We reuse generateStageLinePatchString because the patch structure is the same.
		_, writeErr = tempPatchFile.WriteString(generateDiscardLinePatchString(diffContentStringArray, actualDiscardLineIndex))
		if writeErr != nil {
			tempPatchFile.Close() // Close before return
			return
		}
		// Close the file descriptor so Git can read it safely
		closeErr := tempPatchFile.Close()
		if closeErr != nil {
			return
		}

		switch stageStatus {
		case STAGE:
			// Discarding a change that is already STAGED means unstaging it.
			// We apply the reverse patch to the index (--cached).
			gitArgs = []string{"apply", "--cached", "--unidiff-zero", "--recount", "--whitespace=nowarn", tempPatchFile.Name()}
		case UNSTAGE:
			// Discarding a change that is UNSTAGED (only in worktree) means discarding it from the worktree.
			gitArgs = []string{"apply", "--unidiff-zero", "--recount", "--whitespace=nowarn", tempPatchFile.Name()}
		}

		discardLineCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
		discardLineCmdExecutor.Run()
		// Update file status to reflect changes
		gf.GetGitFilesStatus()
		gf.updateChannel <- GIT_EDIT_LINE_DETAILS_AND_FILES_UPDATE

	}
}

// generateDiscardLinePatchString constructs a unified diff patch that inverses the change of a single line.
// To make the patch apply cleanly even if there are other changes in the same chunk:
// 1. The TARGET line's operator is swapped (+ becomes -, - becomes +) to "undo" it.
// 2. OTHER added lines (+) are converted to context space (' ') because they already exist in the target state.
// 3. OTHER deleted lines (-) are skipped because they do not exist in the target state.
func generateDiscardLinePatchString(diffContentStringArray []string, actualDiscardLineIndex int) string {
	var discardLinePatchString strings.Builder

	lastLineWasSkipped := false
	for index, diffLine := range diffContentStringArray {
		if utf8.RuneCountInString(diffLine) > 0 {
			// Handle special markers
			if strings.HasPrefix(diffLine, "\\ No newline") {
				if lastLineWasSkipped {
					continue
				}
				discardLinePatchString.WriteString(diffLine)
				discardLinePatchString.WriteString("\n")
				continue
			}
			// TARGET LINE: Swap the operator to reverse the change
			if index == actualDiscardLineIndex {
				if utf8.RuneCountInString(diffLine) > 1 {
					if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "---") {
						// Undoing a deletion means adding it back
						discardLinePatchString.WriteString("+")
						discardLinePatchString.WriteString(diffLine[1:])
						discardLinePatchString.WriteString("\n")
					} else if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++") {
						// Undoing an addition means removing it
						discardLinePatchString.WriteString("-")
						discardLinePatchString.WriteString(diffLine[1:])
						discardLinePatchString.WriteString("\n")
					} else {
						discardLinePatchString.WriteString(diffLine)
						discardLinePatchString.WriteString("\n")
					}
				} else {
					// Handle empty lines with just the +/- operator
					if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "---") {
						discardLinePatchString.WriteString("+\n")
					} else if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++") {
						discardLinePatchString.WriteString("-\n")
					} else {
						discardLinePatchString.WriteString(" \n")
					}
				}
				lastLineWasSkipped = false
			} else {
				// OTHER LINES in the chunk: We "normalize" them so the patch only targets our line.
				if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++") {
					// Convert addition to context: this line is already in the file/index, so
					// for the patch to apply, it should be treated as existing context.
					if utf8.RuneCountInString(diffLine) > 1 {
						discardLinePatchString.WriteString(" ")
						discardLinePatchString.WriteString(diffLine[1:])
						discardLinePatchString.WriteString("\n")
					} else {
						discardLinePatchString.WriteString(" \n")
					}
					lastLineWasSkipped = false
				} else if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "---") {
					// Skip deletion: this line is already gone in the file/index, so it
					// doesn't exist to be used as context or modified.
					lastLineWasSkipped = true
					continue
				} else {
					// Keep existing context lines as is.
					discardLinePatchString.WriteString(diffLine)
					discardLinePatchString.WriteString("\n")
					lastLineWasSkipped = false
				}
			}
		} else {
			discardLinePatchString.WriteString("\n")
			lastLineWasSkipped = false
		}
	}

	return discardLinePatchString.String()
}
