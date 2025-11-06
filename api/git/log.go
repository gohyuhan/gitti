package git

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"gitti/cmd"

	"github.com/google/uuid"
)

var GITCOMMITLOG *GitCommitLog

const (
	mergePipe     = "╲"
	verticalPipe  = "│"
	commitNode    = "●"
	branchOutPipe = "╱"
)

// CommitLogInfo holds essential details for a single Git commit.
type CommitLogInfo struct {
	CommitGraphLines []string
	ShortHash        string   // Unique SHA-1 hash of the commit.
	Hash             string   // FUll Unique SHA-1 hash of the commit.
	ParentHashes     []string // Hashes of the parent commits.
	Branches         []string // Local and remote branches that point to this commit (including HEAD).
	Tags             []string // Tags that point to this commit.
	Author           string   // The author's name.
	Message          string   // The commit subject line.
	DateTime         string   // The committer date in ISO 8601 format.
}

type GitCommitLog struct {
	ID                          uuid.UUID // because this is a slow operation and to prevent race override, we use ID to detect and only allow override whe the ID matches
	GitCommitLogs               []CommitLogInfo
	GitCommitLogsDAG            [][]string
	ErrorLog                    []error
	GetCurrentBranchOrAllCommit string
	AllBranch                   bool
	MU                          sync.Mutex
	Status                      string
}

const (
	CURRENTBRANCH = "CURRENT"
	ALLBRANCH     = "ALL"
)

const (
	INITIALIZING = "INITIALIZING"
	DONE         = "DONE"
)

func InitGitCommitLog(allBranch bool) {
	branch := CURRENTBRANCH
	if allBranch {
		branch = ALLBRANCH
	}
	gitCommitLog := GitCommitLog{
		ID:                          uuid.New(),
		GitCommitLogs:               make([]CommitLogInfo, 0, 1),
		GitCommitLogsDAG:            make([][]string, 0, 1),
		ErrorLog:                    []error{},
		GetCurrentBranchOrAllCommit: branch,
		AllBranch:                   allBranch,
	}

	GITCOMMITLOG = &gitCommitLog
}

func (gc *GitCommitLog) GetLatestGitCommitLogInfoAndDAG(updateChannel chan string) {
	gc.ID = uuid.New()
	gc.MU.Lock()
	gc.Status = INITIALIZING
	currentID := gc.ID
	gc.GitCommitLogsDAG = make([][]string, 0, 1)
	gc.GitCommitLogs = make([]CommitLogInfo, 0, 1)
	gc.getCommitLogInfo(currentID, updateChannel)
	gc.Status = DONE
	gc.MU.Unlock()
}

func (gc *GitCommitLog) getCommitLogInfo(currentID uuid.UUID, updateChannel chan string) {
	// Define the custom separators for fields and records.
	const fieldSeparator = "<<COMMITLOGINFOSEPERATOR>>"
	const recordSeparator = "<<COMMITLOGENDOFLINE>>"

	// early return is Id is not the same when the process starts
	if gc.ID != currentID {
		return
	}

	// A custom format with unique separators to handle multi-line commit messages reliably.
	// Format: <hash>|||<parents>|||<refs>|||<author>|||<subject>|||<body>|||<datetime><<RECORD_END>>
	const gitLogFormat = "--pretty=format:" + fieldSeparator + "%h" + fieldSeparator + "%H" + fieldSeparator + "%P" + fieldSeparator + "%d" + fieldSeparator + "%an" + fieldSeparator + "%s" + fieldSeparator + "%ci" + recordSeparator

	// Construct the `git log` command.
	gitArgs := []string{"--no-pager", "log", "--graph", "--topo-order", "--color=always", gitLogFormat}
	if gc.GetCurrentBranchOrAllCommit == ALLBRANCH {
		gitArgs = append(gitArgs, "--all")
	} else {
		gitArgs = append(gitArgs, "HEAD")
	}

	cmd := cmd.GittiCmd.RunGitCmd(gitArgs)
	gitStreamOutput, err := cmd.StdoutPipe()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT LOG ERROR]: %w", err))
		return
	}
	cmd.Stderr = cmd.Stdout

	// to get line back from stdin out streaming
	scanner := bufio.NewScanner(gitStreamOutput)

	// preallocate space
	gc.GitCommitLogs = make([]CommitLogInfo, 0)
	var pendingGraphLines []string

	//start the streaming
	err = cmd.Start()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT LOG ERROR]: %w", err))
		return
	}

	for scanner.Scan() {
		if gc.ID != currentID {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		pendingGraphLines = append(pendingGraphLines, line)
		if strings.Contains(line, recordSeparator) {
			record := line
			record = strings.TrimSuffix(strings.TrimSpace(record), recordSeparator)
			// Split the record into its constituent fields.
			logParts := strings.Split(record, fieldSeparator)
			if len(logParts) < 8 {
				gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT LOG ERROR]: invalid git log format"))
				return
			}

			commitGraphLines := []string{}
			// part 0: the node visualization
			commitNodeGraphLine := strings.TrimSpace(logParts[0])
			// check if there are any more "branch pipe" for between this and the it child commit
			if len(pendingGraphLines) > 1 {
				commitGraphLines = append(pendingGraphLines[0:len(pendingGraphLines)-1], commitNodeGraphLine)
			} else {
				commitGraphLines = append(commitGraphLines, commitNodeGraphLine)
			}
			for i, graphLine := range commitGraphLines {
				commitGraphLines[i] = strings.ReplaceAll(graphLine, "*", commitNode)
				commitGraphLines[i] = strings.ReplaceAll(commitGraphLines[i], "\\", mergePipe)
				commitGraphLines[i] = strings.ReplaceAll(commitGraphLines[i], "/", branchOutPipe)
				commitGraphLines[i] = strings.ReplaceAll(commitGraphLines[i], "|", verticalPipe)

			}
			pendingGraphLines = []string{}

			// Part 1: Parse the short commit hash.
			shortCommitHash := strings.TrimSpace(logParts[1])

			// Part 2: Parse the full commit hash.
			fullCommitHash := strings.TrimSpace(logParts[2])

			// Part 3: Parse parent hashes.
			parentHashesStr := strings.TrimSpace(logParts[3])
			var parentHashes []string
			if parentHashesStr != "" {
				parentHashes = strings.Split(parentHashesStr, " ")
			}

			// Part 4: Parse references (branches and tags).
			referencesStr := strings.TrimSpace(logParts[4])
			var branches, tags []string
			if referencesStr != "" && referencesStr != "()" {
				referencesStr = strings.Trim(referencesStr, "() ")
				refs := strings.SplitSeq(referencesStr, ",")
				for reference := range refs {
					cleanRef := strings.TrimSpace(reference)
					if strings.HasPrefix(cleanRef, "tag: ") {
						tags = append(tags, strings.TrimPrefix(cleanRef, "tag: "))
					} else {
						// Add the reference directly to preserve "HEAD -> main" etc.
						branches = append(branches, cleanRef)
					}
				}
			}

			// Parts 5-7: Parse other commit details
			author := logParts[5] // Defer trimming to when needed
			message := logParts[6]
			dateTime := logParts[7]
			if author != "" {
				author = strings.TrimSpace(author)
			}
			if message != "" {
				message = strings.TrimSpace(message)
			}
			if dateTime != "" {
				dateTime = strings.TrimSpace(dateTime)
			}

			// Create and append the CommitLogInfo struct.
			gc.GitCommitLogs = append(gc.GitCommitLogs, CommitLogInfo{
				CommitGraphLines: commitGraphLines,
				ShortHash:        shortCommitHash,
				Hash:             fullCommitHash,
				ParentHashes:     parentHashes,
				Branches:         branches,
				Tags:             tags,
				Author:           author,
				Message:          message,
				DateTime:         dateTime,
			})

			updateChannel <- GIT_COMMIT_LOG_UPDATE
		}
	}

	// wait for proper clean up
	err = cmd.Wait()
	if err != nil {
		gc.ErrorLog = append(gc.ErrorLog, fmt.Errorf("[GIT COMMIT LOG ERROR]: %w", err))
	}
}
