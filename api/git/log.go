package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gohyuhan/gitti/executor"
)

// Commit represents a single git commit with all necessary graph information
type CommitLog struct {
	Hash       string
	Parents    []string
	Message    string
	Author     string
	LaneString string
	Color      string
}

type GitCommitLog struct {
	errorLog           []error
	gitCommitLogOutput []CommitLog
	updateChannel      chan string
	gitProcessLock     *GitProcessLock
}

// ----------------------------------
//
//	Init Git Commit Log
//
// ----------------------------------
func InitGitCommitLog(updateChannel chan string, gitProcessLock *GitProcessLock) *GitCommitLog {
	gitCommitLog := GitCommitLog{
		gitCommitLogOutput: make([]CommitLog, 0),
		gitProcessLock:     gitProcessLock,
		updateChannel:      updateChannel,
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
		"--pretty=format:%H|%P|%s|%an",
		"-n", "2500",
	}

	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	// Use pipe to process line-by-line to avoid loading entire history into memory
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT LOG ERROR]: %s", err.Error()))
		return
	}

	if err := cmd.Start(); err != nil {
		gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT LOG ERROR]: %s", err.Error()))
		return
	}

	scanner := bufio.NewScanner(stdout)
	renderer := NewGraphRenderer()
	gitCommitLogOutput := make([]CommitLog, 0)
	// 2. Process Commits
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 4)
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
		commitLaneString, nodeColor := renderer.RenderCommit(cL)

		cL.LaneString = commitLaneString
		cL.Color = nodeColor
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

// Palette provides distinct colors for different branches.
var Palette = []string{
	"\033[38;5;57m",  // Deep Purple
	"\033[38;5;63m",  // Purple-Blue
	"\033[38;5;69m",  // Blue
	"\033[38;5;75m",  // Light Blue
	"\033[38;5;81m",  // Cyan-Blue
	"\033[38;5;87m",  // Cyan
	"\033[38;5;49m",  // Spring Green
	"\033[38;5;42m",  // Green
	"\033[38;5;202m", // OrangeRed
	"\033[38;5;196m", // Red
	"\033[38;5;208m", // Orange
	"\033[38;5;220m", // Gold
}

func getColor(colorID int) string {
	if colorID < 0 {
		return ColorReset
	}
	return Palette[colorID%len(Palette)]
}

// --- Graph Renderer ---

// Cell represents a single character in the terminal output grid.
type Cell struct {
	Char  rune
	Color string
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
func (g *GraphRenderer) RenderCommit(cL CommitLog) (string, string) {
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

	// The Node (*) always takes the color of its own lane.
	nodeColor := getColor(commitLane.ColorID)

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
				break
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
		cells[k] = Cell{Char: ' ', Color: ColorReset}
	}

	// Helper to set a character at a specific visual index
	setChar := func(idx int, r rune, color string) {
		if idx >= 0 && idx < len(cells) {
			cells[idx] = Cell{Char: r, Color: color}
		}
	}

	// Helper to draw horizontal lines '─'
	drawHorizontal := func(srcIdx, destIdx int, color string) {
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
			cells[k] = Cell{Char: '─', Color: color}
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
			break
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

		color := getColor(lane.ColorID)

		if nextIdx == -1 {
			// Should not happen for a pass-through (unless it dies unexpectedly),
			// but fallback to straight pipe.
			setChar(i*2, '│', color)
		} else if nextIdx < i {
			// Shifting Left: ↙
			// Visually points to the column it will occupy on the next line.
			setChar(i*2, '↙', color)
		} else if nextIdx > i {
			// Shifting Right: ↘
			setChar(i*2, '↘', color)
		} else {
			// Straight: │
			setChar(i*2, '│', color)
		}
	}

	// Drawing Layer 2: Incoming Merges (Other lanes joining THIS commit)
	for _, srcIdx := range incomingMergeIndices {
		color := getColor(g.currentLanes[srcIdx].ColorID)

		// Draw Horizontal connection to the Commit Node
		drawHorizontal(srcIdx, commitLaneIdx, color)

		// Draw the Corner
		cornerChar := '┘'
		if srcIdx < commitLaneIdx {
			cornerChar = '└'
		}
		setChar(srcIdx*2, cornerChar, color)
	}

	// Drawing Layer 3: Forks (Commit splitting to new Parents)
	// We only draw explicit connectors for Parent 1..N.
	// Parent 0 is implicit (vertical flow).
	if len(parents) > 1 {
		for i := 1; i < len(parents); i++ {
			destIdx := forkDestinations[i] // Where this parent lands in nextLanes

			color := nodeColor

			// Draw Horizontal connection
			drawHorizontal(commitLaneIdx, destIdx, color)

			// Draw Corner at Destination
			cornerChar := '┐'
			if destIdx < commitLaneIdx {
				cornerChar = '┌'
			}
			setChar(destIdx*2, cornerChar, color)
		}
	}

	// Drawing Layer 4: The Commit Node
	commitNodeIndicator := '●'
	if len(parents) > 1 {
		commitNodeIndicator = '◎' // Bullseye for merges
	}
	setChar(commitLaneIdx*2, commitNodeIndicator, nodeColor)

	// -- Step 4: Finalize Output --
	var sb strings.Builder
	for _, c := range cells {
		if c.Char == ' ' {
			sb.WriteString(" ") // standard space (no color)
		} else {
			sb.WriteString(c.Color)
			sb.WriteRune(c.Char)
			sb.WriteString(ColorReset)
		}
	}

	// Update State for next iteration ("Snap" happens here implicitly)
	g.currentLanes = nextLanes

	return strings.TrimRight(sb.String(), " "), nodeColor
}

func (gCL *GitCommitLog) GitCommitLogDetail(ctx context.Context, commitHash string) []string {
	var gitArgs []string

	if gCL.checkIsLargeCommit(commitHash) {
		gitArgs = []string{"show", commitHash}
	} else {
		gitArgs = []string{"show", "--stat", commitHash}
	}

	cmdExecutor := executor.GittiCmdExecutor.RunGitCmdWithContext(ctx, gitArgs, true)
	gitOutput, err := cmdExecutor.Output()
	if err != nil {
		if ctx.Err() != nil {
			// This catches context.Canceled
			gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[COMMIT LOG DETAIL OPERATION CANCELLED DUE TO CONTEXT SWITCHING]: %w", ctx.Err()))
			return nil
		}
		exitError, ok := err.(*exec.ExitError)
		if ok {
			if exitError.ExitCode() != 1 {
				gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT COMMIT LOG DETAIL ERROR]: %w", err))
				return nil
			}
		} else {
			gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT COMMIT LOG DETAIL ERROR]: %w", err))
			return nil
		}
	}

	commitChangesLine := processGeneralGitOpsOutputIntoStringArray(gitOutput)
	return commitChangesLine
}

// ----------------------------------
//
// # Helper to determine if it was a large commit
//
// ----------------------------------
func (gCL *GitCommitLog) checkIsLargeCommit(commitHash string) bool {
	const fileThreshold = 25

	gitArgs := []string{"show", "--shortstat", "--format=''", commitHash}
	cmd := executor.GittiCmdExecutor.RunGitCmd(gitArgs, false)
	cmdOutput, cmdErr := cmd.Output()

	if cmdErr != nil {
		gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT LOG CHECK LARGE COMMIT ERROR]: %s", cmdErr.Error()))
		return true
	}

	re := regexp.MustCompile(`(\d+)\s+files?\s+changed`)
	matches := re.FindStringSubmatch(string(cmdOutput))
	if len(matches) < 2 {
		// No shortstat (e.g. merge commit with no changes)
		return false
	}

	filesChanged, err := strconv.Atoi(matches[1])
	if err != nil {
		gCL.errorLog = append(gCL.errorLog, fmt.Errorf("[GIT LOG CHECK LARGE COMMIT ERROR]: %s", cmdErr.Error()))
		return true
	}

	return filesChanged > fileThreshold
}
