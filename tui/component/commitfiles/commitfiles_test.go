package commitfiles

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

func TestGitCommitFileItem(t *testing.T) {
	item1 := GitCommitFileItem{
		CommitHash:   "abc1234",
		Status:       "M",
		FilePathname: "pkg/core.go",
	}
	if item1.FilterValue() != "pkg/core.go" {
		t.Errorf("expected FilterValue 'pkg/core.go', got '%s'", item1.FilterValue())
	}

	itemRename := GitCommitFileItem{
		CommitHash:   "abc1234",
		Status:       "R100",
		FilePathname: "new.go",
		OldPathname:  "old.go",
	}
	if itemRename.FilterValue() != "old.go -> new.go" {
		t.Errorf("expected FilterValue 'old.go -> new.go', got '%s'", itemRename.FilterValue())
	}
}

func TestGitCommitFileItemDelegateRender(t *testing.T) {
	delegate := GitCommitFileItemDelegate{}
	if delegate.Height() != 1 {
		t.Errorf("expected Height 1, got %d", delegate.Height())
	}
	if delegate.Spacing() != 0 {
		t.Errorf("expected Spacing 0, got %d", delegate.Spacing())
	}

	statuses := []struct {
		status string
		path   string
		old    string
	}{
		{"A", "added.go", ""},
		{"M", "modified.go", ""},
		{"D", "deleted.go", ""},
		{"R100", "renamed.go", "original.go"},
		{"?", "untracked.go", ""},
	}

	var items []list.Item
	for _, s := range statuses {
		items = append(items, GitCommitFileItem{
			CommitHash:   "1234567",
			Status:       s.status,
			FilePathname: s.path,
			OldPathname:  s.old,
		})
	}

	l := list.New(items, delegate, 80, 20)

	for i, item := range items {
		var buf bytes.Buffer
		delegate.Render(&buf, l, i, item)
		rendered := buf.String()
		if len(rendered) == 0 {
			t.Errorf("expected rendered content for index %d, got empty", i)
		}
	}
}

func TestInitCommitFilesList(t *testing.T) {
	dir := t.TempDir()
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

	run("init", "-b", "main")
	run("config", "user.name", "Test")
	run("config", "user.email", "test@example.com")

	os.WriteFile(filepath.Join(dir, "alpha.go"), []byte("alpha"), 0o644)
	os.WriteFile(filepath.Join(dir, "beta.go"), []byte("beta"), 0o644)
	run("add", "alpha.go", "beta.go")
	run("commit", "-m", "add alpha and beta")
	hash := run("rev-parse", "HEAD")

	executor.InitCmdExecutor(dir)
	i18n.InitGittiLanguageMapping("en")
	logger := logging.InitGittiLogging(100, make(chan string, 10), 10)
	lock := git.InitGitProcessLock(logger)
	commitLog := git.InitGitCommitLog(make(chan string, 10), lock, 100, logger)

	m := &types.GittiModel{
		RepoPath: dir,
		GitOperations: &api.GitOperations{
			GitCommitLog: commitLog,
		},
		WindowLeftPanelWidth:          60,
		CommitLogComponentPanelHeight: 20,
		PanelFilterQuery:              make(map[string]string),
	}

	ctx := context.Background()
	hasFiles := InitCommitFilesList(ctx, m, hash, "add alpha and beta")
	if !hasFiles {
		t.Fatalf("expected InitCommitFilesList to return true")
	}
	if len(m.CurrentRepoCommitFilesList.Items()) != 2 {
		t.Errorf("expected 2 items, got %d", len(m.CurrentRepoCommitFilesList.Items()))
	}

	// Test with filter query active
	m.PanelFilterQuery[constant.CommitFilesComponentPanel] = "alpha"
	InitCommitFilesList(ctx, m, hash, "add alpha and beta")
	if len(m.CurrentRepoCommitFilesList.Items()) != 1 {
		t.Errorf("expected 1 item after filtering, got %d", len(m.CurrentRepoCommitFilesList.Items()))
	}
	filteredItem, ok := m.CurrentRepoCommitFilesList.SelectedItem().(GitCommitFileItem)
	if !ok || filteredItem.FilePathname != "alpha.go" {
		t.Errorf("expected selected item 'alpha.go', got %+v", filteredItem)
	}
}
