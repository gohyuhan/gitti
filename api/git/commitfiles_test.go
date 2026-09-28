package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
)

func setupTestGitRepo(t *testing.T) (string, *GitCommitLog) {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s failed: %v, output: %s", strings.Join(args, " "), err, string(out))
		}
		return strings.TrimSpace(string(out))
	}

	run("init", "-b", "main")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@example.com")

	i18n.InitGittiLanguageMapping("en")
	executor.InitCmdExecutor(dir)
	logger := logging.InitGittiLogging(100, make(chan string, 10), 10)
	lock := InitGitProcessLock(logger)
	commitLog := InitGitCommitLog(make(chan string, 10), lock, 100, logger)

	return dir, commitLog
}

func TestCommitTouchedFilesAndOperations(t *testing.T) {
	dir, commitLog := setupTestGitRepo(t)
	ctx := context.Background()

	run := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s failed: %v, output: %s", strings.Join(args, " "), err, string(out))
		}
		return strings.TrimSpace(string(out))
	}

	// 1. Initial Root Commit
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("initial file 1\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("initial file 2\n"), 0o644)
	run("add", "file1.txt", "file2.txt")
	run("commit", "-m", "Initial root commit")
	rootHash := run("rev-parse", "HEAD")

	rootFiles := commitLog.GetCommitTouchedFiles(ctx, rootHash)
	if len(rootFiles) != 2 {
		t.Fatalf("expected 2 files in root commit, got %d", len(rootFiles))
	}
	for _, f := range rootFiles {
		if f.Status != "A" {
			t.Errorf("expected status 'A' for root commit file, got %s for %s", f.Status, f.FilePathname)
		}
	}

	// 2. Second Commit: modify file1, delete file2, add file3
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("modified file 1\n"), 0o644)
	run("rm", "file2.txt")
	os.WriteFile(filepath.Join(dir, "file3.txt"), []byte("new file 3\n"), 0o644)
	run("add", "file1.txt", "file3.txt")
	run("commit", "-m", "Second commit: mod, del, add")
	secondHash := run("rev-parse", "HEAD")

	secondFiles := commitLog.GetCommitTouchedFiles(ctx, secondHash)
	if len(secondFiles) != 3 {
		t.Fatalf("expected 3 files in second commit, got %d", len(secondFiles))
	}

	statusMap := make(map[string]string)
	for _, f := range secondFiles {
		statusMap[f.FilePathname] = f.Status
	}

	if statusMap["file1.txt"] != "M" {
		t.Errorf("expected file1.txt status 'M', got '%s'", statusMap["file1.txt"])
	}
	if statusMap["file2.txt"] != "D" {
		t.Errorf("expected file2.txt status 'D', got '%s'", statusMap["file2.txt"])
	}
	if statusMap["file3.txt"] != "A" {
		t.Errorf("expected file3.txt status 'A', got '%s'", statusMap["file3.txt"])
	}

	// 3. Test GetCommitFileDiff
	diffOutput := commitLog.GetCommitFileDiff(ctx, secondHash, "file1.txt")
	diffJoined := strings.Join(diffOutput, "\n")
	if !strings.Contains(diffJoined, "modified file 1") {
		t.Errorf("expected diff to contain 'modified file 1', got:\n%s", diffJoined)
	}

	// 4. Test CheckoutFileFromCommit:
	// Working directory has "modified file 1". Check out file1 from root commit.
	err := commitLog.CheckoutFileFromCommit(ctx, rootHash, "file1.txt")
	if err != nil {
		t.Fatalf("CheckoutFileFromCommit failed: %v", err)
	}
	content, _ := os.ReadFile(filepath.Join(dir, "file1.txt"))
	if string(content) != "initial file 1\n" {
		t.Errorf("expected file1.txt to have root commit content after checkout, got %s", string(content))
	}

	// Revert working tree back to second commit state
	run("checkout", "HEAD", "--", "file1.txt")

	// 5. Test DiscardFileFromCommit:
	// Discarding file1.txt from second commit should restore it to root commit state (HEAD~1)
	err = commitLog.DiscardFileFromCommit(ctx, secondHash, "file1.txt")
	if err != nil {
		t.Fatalf("DiscardFileFromCommit failed: %v", err)
	}
	content, _ = os.ReadFile(filepath.Join(dir, "file1.txt"))
	if string(content) != "initial file 1\n" {
		t.Errorf("expected file1.txt to have root commit content after discard, got %s", string(content))
	}

	// 5b. Test DiscardFileFromCommit for newly added file in non-root commit:
	// Discarding file3.txt (added in second commit) should remove it from working directory
	err = commitLog.DiscardFileFromCommit(ctx, secondHash, "file3.txt")
	if err != nil {
		t.Fatalf("DiscardFileFromCommit on added file in second commit failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "file3.txt")); !os.IsNotExist(err) {
		t.Errorf("expected file3.txt to be removed after discard, but it exists")
	}

	// 6. Test DiscardFileFromCommit on root commit:
	// Discarding a file added in root commit should remove the file from working directory
	err = commitLog.DiscardFileFromCommit(ctx, rootHash, "file1.txt")
	if err != nil {
		t.Fatalf("DiscardFileFromCommit on root commit failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "file1.txt")); !os.IsNotExist(err) {
		t.Errorf("expected file1.txt to be removed after root commit discard, but it exists")
	}
}

func TestCommitTouchedFilesRename(t *testing.T) {
	dir, commitLog := setupTestGitRepo(t)
	ctx := context.Background()

	run := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s failed: %v, output: %s", strings.Join(args, " "), err, string(out))
		}
		return strings.TrimSpace(string(out))
	}

	os.WriteFile(filepath.Join(dir, "old_name.txt"), []byte("content for rename test\n"), 0o644)
	run("add", "old_name.txt")
	run("commit", "-m", "Commit before rename")

	run("mv", "old_name.txt", "new_name.txt")
	run("commit", "-m", "Rename file")
	renameCommitHash := run("rev-parse", "HEAD")

	files := commitLog.GetCommitTouchedFiles(ctx, renameCommitHash)
	if len(files) != 1 {
		t.Fatalf("expected 1 touched file for rename, got %d", len(files))
	}

	if !strings.HasPrefix(files[0].Status, "R") {
		t.Errorf("expected status starting with 'R', got '%s'", files[0].Status)
	}
	if files[0].FilePathname != "new_name.txt" {
		t.Errorf("expected FilePathname 'new_name.txt', got '%s'", files[0].FilePathname)
	}
	if files[0].OldPathname != "old_name.txt" {
		t.Errorf("expected OldPathname 'old_name.txt', got '%s'", files[0].OldPathname)
	}
}
