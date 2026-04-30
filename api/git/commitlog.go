package git

import (
	"bufio"
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
)

const SEPARATOR = "\x00"

// represent the info of commit that got cherry picked
type CherryPickedCommitLog struct {
	Hash                 string
	Message              string
	Author               string
	FromBranch           string
	UserSelectedSequence int
}

// Commit represents a single git commit with all necessary graph information
type CommitLog struct {
	Hash         string
	Parents      []string
	Message      string
	Author       string
	LaneCharInfo []Cell
	ColorID      int
}

type GitCommitLog struct {
	gitCommitLogOutput []CommitLog
	updateChannel      chan string
	maxCommitLogCount  string
	gitProcessLock     *GitProcessLock
	logging            *logging.GittiLogging
}

type CommitHashParentInfo struct {
	ParentCommitHash    string
	ParentCommitMessage string
	ParentOrder         int
}

// ----------------------------------
//
//	Init Git Commit Log
//
// ----------------------------------
func InitGitCommitLog(updateChannel chan string, gitProcessLock *GitProcessLock, maxCommitLogCountInt int, logging *logging.GittiLogging) *GitCommitLog {
	maxCommitLogCount := strconv.Itoa(maxCommitLogCountInt)
	gitCommitLog := GitCommitLog{
		gitCommitLogOutput: make([]CommitLog, 0),
		gitProcessLock:     gitProcessLock,
		updateChannel:      updateChannel,
		maxCommitLogCount:  maxCommitLogCount,
		logging:            logging,
	}
	return &gitCommitLog
}

// ----------------------------------
//
//	Return commit log output
//
// ----------------------------------
func (gCL *GitCommitLog) GitCommitLogOutput() []CommitLog {
	copied := make([]CommitLog, len(gCL.gitCommitLogOutput))
	copy(copied, gCL.gitCommitLogOutput)
	return copied
}

// ----------------------------------
//
//	Get the Commit log
//
// ----------------------------------
func (gCL *GitCommitLog) GetCommitLogs() {
	// 1. Prepare git command
	gitArgs := []string{
		"log",
		"--topo-order",
		"--no-decorate",
		"--no-notes",
		"--pretty=format:%H%x00%P%x00%s%x00%an",
		"-n", gCL.maxCommitLogCount,
		"--",
	}

	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	// Use pipe to process line-by-line to avoid loading entire history into memory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gCL.logging.RegisterNewLog(logging.COMMIT_LOG_OPS, "", logging.ERROR, fmt.Sprintf("[PIPE ERROR]: %s", err.Error()), false)
		return
	}

	if err := cmd.Start(); err != nil {
		gCL.logging.RegisterNewLog(logging.COMMIT_LOG_OPS, "", logging.ERROR, fmt.Sprintf("[START ERROR]: %s", err.Error()), false)
		return
	}

	scanner := bufio.NewScanner(stdout)
	renderer := NewGraphRenderer()
	gitCommitLogOutput := make([]CommitLog, 0)
	// 2. Process Commits
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, SEPARATOR, 4)
		if len(parts) < 4 {
			continue
		}

		// Parse commit
		cL := CommitLog{
			Hash:    parts[0],
			Message: parts[2],
			Author:  parts[3],
		}
		if len(parts[1]) > 0 {
			cL.Parents = strings.Split(parts[1], " ")
		}

		// 3. Render
		// The renderer returns the commit lane string
		laneCharInfo, colorID := renderer.RenderCommit(cL)

		cL.LaneCharInfo = laneCharInfo
		cL.ColorID = colorID
		gitCommitLogOutput = append(gitCommitLogOutput, cL)
	}

	gCL.gitCommitLogOutput = gitCommitLogOutput
}

// RenderCommit generates the visual graph line for a single commit.
//
// Algorithm Overview: "Stable-Color Dense-Packing"
// ------------------------------------------------
//  1. One Line Per Commit: We draw exactly one line of text for this commit.
//  2. Dense Packing: We do not leave holes in the lane list. If a branch dies (merges),
//     branches to its right "snap" to the left immediately on the next line.
//  3. Stable Colors: To make this "snap" less confusing, branches keep their assigned
//     color even if they move to a different column index.

// --- Graph Renderer ---

// Cell represents a single character in the terminal output grid.
type Cell struct {
	Char    rune
	ColorID int
}

// Lane represents a persistent branch in the graph.
// Important: The ColorID serves as the unique identifier for visual continuity.
type Lane struct {
	Hash    string
	ColorID int
}

type GraphRenderer struct {
	// currentLanes tracks the state of branches at the CURRENT print line.
	// It is "Dense", meaning there are no empty gaps/nil entries.
	currentLanes []Lane
}

// ----------------------------------
//
//	Create a new GraphRenderer for rendering commit lane visualization
//
// ----------------------------------
func NewGraphRenderer() *GraphRenderer {
	return &GraphRenderer{
		currentLanes: make([]Lane, 0),
	}
}

// ----------------------------------
//
//	Render the Commit Lane graph line
//
// ----------------------------------
func (g *GraphRenderer) RenderCommit(cL CommitLog) ([]Cell, int) {
	// -- Step 1: Identify the Commit's Lane --
	// Find which existing lane this commit belongs to.

	commitLaneIdx := -1
	var commitLane Lane

	for i, l := range g.currentLanes {
		if l.Hash == cL.Hash {
			commitLaneIdx = i
			commitLane = l
			break
		}
	}

	// Case: New Tip (Root or independent branch start)
	// If the commit isn't in our tracked lanes, it's a new starting point.
	if commitLaneIdx == -1 {
		// Create a new Lane identity.
		commitLane = Lane{
			Hash:    cL.Hash,
			ColorID: len(g.currentLanes),
		}

		// Append to the rightmost side (Visual preference).
		commitLaneIdx = len(g.currentLanes)
		g.currentLanes = append(g.currentLanes, commitLane)
	}

	// -- Step 2: Plan the Next State (Evolution) --
	// We determine what the lanes will look like for the NEXT commit.
	// This involves:
	// 1. Removing lanes that merge into this commit (they die here). (branching out)
	// 2. Updating the current lane to point to its First Parent.
	// 3. Creating NEW lanes for any additional parents (Forks). (merging in)

	// Identify "Incoming Merges"
	// These are OTHER lanes that point to THIS commit. They will be drawn joining in.
	var incomingMergeIndices []int
	for i, l := range g.currentLanes {
		if i != commitLaneIdx && l.Hash == cL.Hash {
			incomingMergeIndices = append(incomingMergeIndices, i)
		}
	}

	// Build 'nextLanes' (The state for the next iteration).
	// We rebuild this list from scratch to ensure it remains Dense (no gaps).
	var nextLanes []Lane

	// Track where forks need to connect to.
	// Map: ParentIndex (0, 1..) -> Destination Column Index in nextLanes
	forkDestinations := make(map[int]int)

	parents := cL.Parents

	// Iterate through CURRENT lanes to decide their fate.
	for i, l := range g.currentLanes {
		if i == commitLaneIdx {
			// This is the Active Lane for the current commit.
			if len(parents) > 0 {
				// Continuation: The lane continues to Parent 0.
				// It keeps the same ColorID.
				p0Lane := Lane{
					Hash:    parents[0],
					ColorID: commitLane.ColorID,
				}

				// Add to next state
				newIdx := len(nextLanes)
				nextLanes = append(nextLanes, p0Lane)

				// Parent 0 is the "Straight" continuation
				forkDestinations[0] = newIdx
			} else {
				// No parents (Root Commit of repo).
				// The lane ends here. We do NOT add it to nextLanes.
			}
		} else {
			// Check if this lane is merging INTO us.
			isMerge := false
			if slices.Contains(incomingMergeIndices, i) {
				isMerge = true
			}

			if isMerge {
				// It merges here. It dies.
				// clearly visually indicated by a '┘' or '└' connector later.
				// Do NOT add to nextLanes.
			} else {
				// Independent lane (Pass-Through).
				// It just carries over to the next state, keeping its ColorID.
				nextLanes = append(nextLanes, l)
			}
		}
	}

	// Handle Forks (Parents 1..N)
	// These are new branches splitting off from this commit.
	if len(parents) > 1 {
		for pIn := 1; pIn < len(parents); pIn++ {
			pHash := parents[pIn]

			// Start a NEW Lane with a NEW ColorID
			newLane := Lane{
				Hash:    pHash,
				ColorID: len(nextLanes),
			}

			// Append to the list
			newIdx := len(nextLanes)
			nextLanes = append(nextLanes, newLane)
			forkDestinations[pIn] = newIdx
		}
	}

	// -- Step 3: Draw the Current Line --
	// We render the visual connections based on the CURRENT state indices.
	// Use 2 characters per lane width: "| " or "* " etc.

	// Calculate Grid Width
	// We need enough space to draw the current lanes AND any connectors to new forks.
	// Since we only simply append forks, the max width is determined by `nextLanes`.
	maxWidth := max(len(nextLanes), len(g.currentLanes))

	lineLen := maxWidth * 2
	cells := make([]Cell, lineLen+1) // +1 buffer
	// Initialize with empty
	for k := range cells {
		cells[k] = Cell{Char: ' ', ColorID: -1}
	}

	// Helper to set a character at a specific visual index
	setChar := func(idx int, r rune, colorID int) {
		if idx >= 0 && idx < len(cells) {
			cells[idx] = Cell{Char: r, ColorID: colorID}
		}
	}

	// Helper to draw horizontal lines '─'
	drawHorizontal := func(srcIdx, destIdx int, colorID int) {
		// Convert logical indices directly to visual indices (x2)
		start := srcIdx * 2
		end := destIdx * 2

		// Ensure Start < End for loop
		if start > end {
			start, end = end, start
			start += 1 // Adjust bounds to not overwrite the corner characters
			end -= 1
		} else {
			start += 1
			end -= 1
		}

		for k := start; k <= end; k++ {
			// Protection: Don't overwrite any existing character (Pipes, Diagonals, Nodes)
			if cells[k].Char != ' ' {
				continue
			}
			cells[k] = Cell{Char: '─', ColorID: colorID}
		}
	}

	// Drawing Layer 1: Vertical Pipes (Pass-Throughs)
	// These are lanes that are NOT the current commit and NOT merging in.
	checkedNextLanesIndices := make(map[int]bool)
	for i := range g.currentLanes {
		if i == commitLaneIdx {
			continue
		} // Skip active lane (it gets a Node *)

		// Check if it's a merge source
		isMerge := false
		if slices.Contains(incomingMergeIndices, i) {
			isMerge = true
		}

		if isMerge {
			continue
		} // Handled in Layer 2

		// It is a Pass-Through lane.
		// Determine which character to draw based on its Next Position.
		lane := g.currentLanes[i]

		// Find where this lane goes in nextLanes
		nextIdx := -1
		for j, nl := range nextLanes {
			if nl.ColorID == lane.ColorID && !checkedNextLanesIndices[j] {
				checkedNextLanesIndices[j] = true
				nextIdx = j
				break
			}
		}

		if nextIdx == -1 {
			// Should not happen for a pass-through (unless it dies unexpectedly),
			// but fallback to straight pipe.
			setChar(i*2, '│', lane.ColorID)
		} else if nextIdx < i {
			// Shifting Left: ↙
			// Visually points to the column it will occupy on the next line.
			setChar(i*2, '↙', lane.ColorID)
		} else if nextIdx > i {
			// Shifting Right: ↘
			setChar(i*2, '↘', lane.ColorID)
		} else {
			// Straight: │
			setChar(i*2, '│', lane.ColorID)
		}
	}

	// Drawing Layer 2: Incoming Merges (Other lanes joining THIS commit)
	for _, srcIdx := range incomingMergeIndices {
		// Draw Horizontal connection to the Commit Node
		drawHorizontal(srcIdx, commitLaneIdx, g.currentLanes[srcIdx].ColorID)

		// Draw the Corner
		cornerChar := '┘'
		if srcIdx < commitLaneIdx {
			cornerChar = '└'
		}
		setChar(srcIdx*2, cornerChar, g.currentLanes[srcIdx].ColorID)
	}

	// Drawing Layer 3: Forks (Commit splitting to new Parents)
	// We only draw explicit connectors for Parent 1..N.
	// Parent 0 is implicit (vertical flow).
	if len(parents) > 1 {
		for i := 1; i < len(parents); i++ {
			destIdx := forkDestinations[i] // Where this parent lands in nextLanes

			// Draw Horizontal connection
			drawHorizontal(commitLaneIdx, destIdx, commitLaneIdx)

			// Draw Corner at Destination
			cornerChar := '┐'
			if destIdx < commitLaneIdx {
				cornerChar = '┌'
			}
			setChar(destIdx*2, cornerChar, commitLaneIdx)
		}
	}

	// Drawing Layer 4: The Commit Node
	commitNodeIndicator := '●'
	if len(parents) > 1 {
		commitNodeIndicator = '◎' // Bullseye for merges
	}
	setChar(commitLaneIdx*2, commitNodeIndicator, commitLaneIdx)

	// Update State for next iteration ("Snap" happens here implicitly)
	g.currentLanes = nextLanes

	return cells, commitLaneIdx
}

// ----------------------------------
//
//	Retrieve the detailed diff or stat for a specific commit by its hash
//
// ----------------------------------
func (gCL *GitCommitLog) GitCommitLogDetail(ctx context.Context, commitHash string) []string {
	gitArgs := []string{"show", commitHash, "--stat", "-p"}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	gitOutput, err := cmdExecutor.Output()
	if err != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gCL.logging.RegisterNewLog(logging.COMMIT_LOG_DETAIL_OPS, strings.Join(gitArgs, " "), logging.WARN, fmt.Sprintf("[%s CANCELLED]: %s", logging.COMMIT_LOG_DETAIL_OPS, ctx.Err().Error()), true)
			return nil
		} else {
			gCL.logging.RegisterNewLog(logging.COMMIT_LOG_DETAIL_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.COMMIT_LOG_DETAIL_OPS, err.Error()), true)
			return nil
		}
	}

	commitChangesLine := processGeneralGitOpsOutputIntoStringArray(gitOutput)
	return commitChangesLine
}

// ----------------------------------
//
// # Commit cherry pick
//
// ----------------------------------
func (gCL *GitCommitLog) GitCherryPick(cherryPickedCommitHashes []string) {
	topoOrderedCherryPickedCommitHashes := gCL.topoOrderCherryPickedCommit(cherryPickedCommitHashes)
	gitArgs := []string{"cherry-pick"}
	gitArgs = append(gitArgs, topoOrderedCherryPickedCommitHashes...)

	cherryPickCmdExec := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	err := cherryPickCmdExec.Run()
	gCL.logging.RegisterNewLog(logging.CHERRY_PICK_COMMIT_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)
	if err != nil {
		gCL.logging.RegisterNewLog(logging.CHERRY_PICK_COMMIT_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.CHERRY_PICK_COMMIT_OPS, err.Error()), true)
	}
}

// ----------------------------------
//
//	GitCherryPickWithSigning constructs the git cherry-pick command arguments for execution in the terminal.
//	This allows for interactive signing (e.g., GPG passphrase) by suspending the UI.
//
// ----------------------------------
func (gCL *GitCommitLog) GitCherryPickWithSigning(cherryPickedCommitHashes []string) []string {
	topoOrderedCherryPickedCommitHashes := gCL.topoOrderCherryPickedCommit(cherryPickedCommitHashes)
	gitArgs := []string{"cherry-pick"}
	gitArgs = append(gitArgs, topoOrderedCherryPickedCommitHashes...)
	return gitArgs
}

// ----------------------------------
//
// Helper to topo order cherry picked commit to prevent cherry pick conflict
// (this only to commits that are from the same branch or related, else it will be in the sequence of how user cherry oicked it)
//
// ----------------------------------
func (gCL *GitCommitLog) topoOrderCherryPickedCommit(cherryPickedCommitHashes []string) []string {
	// if the cherry picked commit hash is less than 1, we don't have to even order it
	if len(cherryPickedCommitHashes) <= 1 {
		return cherryPickedCommitHashes
	}

	var topoOrderedCherryPickedCommitHashes []string
	gitArgs := []string{"rev-list", "--topo-order", "--reverse", "--no-walk"}
	gitArgs = append(gitArgs, cherryPickedCommitHashes...)

	topoOrderCherryPickedCommitHashesCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	topoOrderCherryPiCkedCommitHashesOutput, topoOrderCherryPiCkedCommitHashesErr := topoOrderCherryPickedCommitHashesCmdExecutor.Output()
	if topoOrderCherryPiCkedCommitHashesErr != nil {
		// will default to current cherry picked commit hashes sequencees by user
		return cherryPickedCommitHashes
	}

	topoOrderedCherryPickedCommitHashes = processGeneralGitOpsOutputIntoStringArray(topoOrderCherryPiCkedCommitHashesOutput)

	return topoOrderedCherryPickedCommitHashes
}

// ----------------------------------
//
// # Helper to write commit graph that improve performance of git log retrieval
//
// ----------------------------------
func (gCL *GitCommitLog) WriteCommitGraph() {
	gitArgs := []string{"commit-graph", "write", "--reachable", "--split"}
	writeCommitGraphCmdExec := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	writeCommitGraphCmdExec.Run()
}

// ----------------------------------
//
// # Git Revert for commit
//
// ----------------------------------
func (gCL *GitCommitLog) GitRevertCommit(commitHash string, parentOrder int) {
	// parentOrder is only need when reverting a merge commit, other commit reverting only require the hash commit
	var gitArgs []string
	if parentOrder > 0 {
		parentOrderString := strconv.Itoa(parentOrder)
		gitArgs = []string{"revert", "--no-edit", "-m", parentOrderString, commitHash}
	} else {
		gitArgs = []string{"revert", "--no-edit", commitHash}
	}

	revertCommitCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	gCL.logging.RegisterNewLog(logging.REVERT_COMMIT_OPS, strings.Join(gitArgs, " "), logging.INFO, "", true)

	err := revertCommitCmdExecutor.Run()
	if err != nil {
		gCL.logging.RegisterNewLog(logging.REVERT_COMMIT_OPS, strings.Join(gitArgs, " "), logging.ERROR, fmt.Sprintf("[%s ERROR]: %s", logging.REVERT_COMMIT_OPS, err.Error()), true)
	}
}

// ----------------------------------
//
// # To get the parent(s) info for the commit hash
//
// ----------------------------------
func (gCL *GitCommitLog) GetCommitHashParentInfo(commitHash string) []CommitHashParentInfo {
	parentTarget := commitHash + "^@"

	gitArgs := []string{"show", "-s", "--format=%H||%s", parentTarget}
	commitHashParentInfoCmdExecutor := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	output, err := commitHashParentInfoCmdExecutor.Output()
	if err != nil {
		return nil
	}

	parsedCommitHashParentInfo := processGeneralGitOpsOutputIntoStringArray(output)

	var parentCommitInfoArray []CommitHashParentInfo
	for index := range parsedCommitHashParentInfo {
		parentCommitInfo := strings.Split(parsedCommitHashParentInfo[index], "||")
		commitHashParentInfo := CommitHashParentInfo{
			ParentCommitHash:    parentCommitInfo[0],
			ParentCommitMessage: parentCommitInfo[1],
			ParentOrder:         index + 1,
		}

		parentCommitInfoArray = append(parentCommitInfoArray, commitHashParentInfo)
	}

	return parentCommitInfoArray
}

// GitRevertCommitWithSigning constructs the git revert command arguments for execution in the terminal.
// This allows for interactive signing (e.g., GPG passphrase) by suspending the UI.
func (gCL *GitCommitLog) GitRevertCommitWithSigning(commitHash string, parentOrder int) []string {
	var gitArgs []string
	if parentOrder > 0 {
		parentOrderString := strconv.Itoa(parentOrder)
		gitArgs = []string{"revert", "--no-edit", "-m", parentOrderString, commitHash}
	} else {
		gitArgs = []string{"revert", "--no-edit", commitHash}
	}
	return gitArgs
}
