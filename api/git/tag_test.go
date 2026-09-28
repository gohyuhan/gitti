package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
)

func setupTestGitRepoWithTag(t *testing.T, tagName string) (*GitTag, *GitBranch, string) {
	t.Helper()
	tempDir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = tempDir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s", args, string(out))
		}
	}

	run("init")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@example.com")

	dummyFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(dummyFile, []byte("hello world\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	run("add", "test.txt")
	run("commit", "-m", "initial commit")
	run("tag", tagName)

	i18n.InitGittiLanguageMapping("EN")
	executor.InitCmdExecutor(tempDir)
	logChan := make(chan string, 100)
	logger := logging.InitGittiLogging(100, logChan, 10)
	lock := InitGitProcessLock(logger)
	gitTag := InitGitTag(logChan, lock, logger)
	gitBranch := InitGitBranch(lock, true, logger)

	return gitTag, gitBranch, tempDir
}

func TestGitCheckoutTag(t *testing.T) {
	gitTag, gitBranch, _ := setupTestGitRepoWithTag(t, "v1.0.0")

	output, success := gitTag.GitCheckoutTag("v1.0.0")
	if !success {
		t.Fatalf("GitCheckoutTag failed, output: %v", output)
	}

	gitBranch.GetLatestBranchesInfo()
	current := gitBranch.CurrentCheckOut()
	if !current.IsCheckedOut {
		t.Errorf("expected branch/head to be checked out")
	}
}

func TestGitCheckoutTag_NonExistent(t *testing.T) {
	gitTag, _, _ := setupTestGitRepoWithTag(t, "v1.0.0")

	output, success := gitTag.GitCheckoutTag("v9.9.9")
	if success {
		t.Fatalf("expected GitCheckoutTag to fail for nonexistent tag, got success with output: %v", output)
	}
}
