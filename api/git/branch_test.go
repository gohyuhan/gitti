package git

import (
	"testing"
)

func TestGitCreateNewBranchBasedOnTagAndSwitch(t *testing.T) {
	_, gitBranch, _ := setupTestGitRepoWithTag(t, "v1.2.3")

	newBranch := "feature-from-tag"
	output, success := gitBranch.GitCreateNewBranchBasedOnTagAndSwitch(newBranch, "v1.2.3")
	if !success {
		t.Fatalf("GitCreateNewBranchBasedOnTagAndSwitch failed, output: %v", output)
	}

	gitBranch.GetLatestBranchesInfo()
	current := gitBranch.CurrentCheckOut()
	if current.BranchName != newBranch {
		t.Errorf("expected current branch to be %s, got %s", newBranch, current.BranchName)
	}
}

func TestGitCreateNewBranchBasedOnTagAndSwitch_NonExistentTag(t *testing.T) {
	_, gitBranch, _ := setupTestGitRepoWithTag(t, "v1.2.3")

	output, success := gitBranch.GitCreateNewBranchBasedOnTagAndSwitch("invalid-branch", "nonexistent-tag")
	if success {
		t.Fatalf("expected failure for nonexistent tag, got success with output: %v", output)
	}
}
