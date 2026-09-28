package nontyping

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/component/commitfiles"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/popup/copypopup"
	"github.com/gohyuhan/gitti/tui/types"
)

func setupTestModel(t *testing.T) (*types.GittiModel, string, string) {
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

	os.WriteFile(filepath.Join(dir, "file_a.txt"), []byte("content a"), 0o644)
	os.WriteFile(filepath.Join(dir, "file_b.txt"), []byte("content b"), 0o644)
	run("add", "file_a.txt", "file_b.txt")
	run("commit", "-m", "initial commit")
	hash := run("rev-parse", "HEAD")

	executor.InitCmdExecutor(dir)
	i18n.InitGittiLanguageMapping("en")
	logger := logging.InitGittiLogging(100, make(chan string, 10), 10)
	lock := git.InitGitProcessLock(logger)
	gitCommitLog := git.InitGitCommitLog(make(chan string, 10), lock, 100, logger)

	commitItems := []list.Item{
		commitlog.GitCommitLogItem{
			Hash:    hash,
			Message: "initial commit",
			Author:  "Test",
		},
	}
	commitList := list.New(commitItems, commitlog.GitCommitLogItemDelegate{}, 60, 20)

	fileItems := []list.Item{
		commitfiles.GitCommitFileItem{
			CommitHash:   hash,
			Status:       "A",
			FilePathname: "file_a.txt",
		},
		commitfiles.GitCommitFileItem{
			CommitHash:   hash,
			Status:       "A",
			FilePathname: "file_b.txt",
		},
	}
	filesList := list.New(fileItems, commitfiles.GitCommitFileItemDelegate{}, 60, 20)

	m := &types.GittiModel{
		RepoPath: dir,
		GitOperations: &api.GitOperations{
			GitCommitLog: gitCommitLog,
		},
		CurrentSelectedComponent:                  constant.CommitLogOrRefLogComponentPanel,
		CurrentCommitLogOrRefLogComponentShowing: constant.SHOW_COMMITLOG,
		CurrentRepoCommitLogInfoList:             commitList,
		CurrentRepoCommitFilesList:               filesList,
		WindowLeftPanelWidth:                     60,
		CommitLogComponentPanelHeight:            20,
		PanelFilterQuery:                         make(map[string]string),
	}
	m.ShowPopUp.Store(false)
	m.IsLineEditingState.Store(false)
	m.IsTyping.Store(false)

	return m, dir, hash
}

func TestCommitFilesDrillDownEnterAndEsc(t *testing.T) {
	m, _, hash := setupTestModel(t)

	// 1. Press Enter on commit log item -> should drill down to CommitFilesComponentPanel
	m, _ = handleNonTypingEnterKeyBindingInteraction(m)
	if m.CurrentSelectedComponent != constant.CommitFilesComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to be CommitFilesComponentPanel, got %s", m.CurrentSelectedComponent)
	}
	if m.CurrentDrillDownCommitHash != hash {
		t.Errorf("expected CurrentDrillDownCommitHash %s, got %s", hash, m.CurrentDrillDownCommitHash)
	}
	if m.CurrentDrillDownCommitSubject != "initial commit" {
		t.Errorf("expected CurrentDrillDownCommitSubject 'initial commit', got %s", m.CurrentDrillDownCommitSubject)
	}

	// 2. Press Enter on commit file item -> should focus DetailComponentPanel with parent CommitFilesComponentPanel
	m, _ = handleNonTypingEnterKeyBindingInteraction(m)
	if m.CurrentSelectedComponent != constant.DetailComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to be DetailComponentPanel, got %s", m.CurrentSelectedComponent)
	}
	if m.DetailPanelParentComponent != constant.CommitFilesComponentPanel {
		t.Fatalf("expected DetailPanelParentComponent to be CommitFilesComponentPanel, got %s", m.DetailPanelParentComponent)
	}

	// 3. Press Esc in DetailComponentPanel -> should return to CommitFilesComponentPanel
	m, _ = handleNonTypingEscKeyBindingInteraction(m)
	if m.CurrentSelectedComponent != constant.CommitFilesComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to return to CommitFilesComponentPanel, got %s", m.CurrentSelectedComponent)
	}

	// 4. Press Esc in CommitFilesComponentPanel -> should return to CommitLogOrRefLogComponentPanel
	m, _ = handleNonTypingEscKeyBindingInteraction(m)
	if m.CurrentSelectedComponent != constant.CommitLogOrRefLogComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to return to CommitLogOrRefLogComponentPanel, got %s", m.CurrentSelectedComponent)
	}
	if m.CurrentDrillDownCommitHash != "" {
		t.Errorf("expected CurrentDrillDownCommitHash to be cleared, got %s", m.CurrentDrillDownCommitHash)
	}
}

func TestCommitFilesNavigation(t *testing.T) {
	m, _, _ := setupTestModel(t)
	m.CurrentSelectedComponent = constant.CommitFilesComponentPanel
	m.CurrentRepoCommitFilesList.Select(0)
	m.ListNavigationIndexPosition.CommitFilesComponent = 0

	// Navigate down
	downMsg := tea.KeyPressMsg{Code: 'j'}
	m, _ = handleNonTypingDownjKeyBindingInteraction(downMsg, m)
	if m.CurrentRepoCommitFilesList.Index() != 1 {
		t.Errorf("expected index 1 after down navigation, got %d", m.CurrentRepoCommitFilesList.Index())
	}
	if m.ListNavigationIndexPosition.CommitFilesComponent != 1 {
		t.Errorf("expected ListNavigationIndexPosition 1, got %d", m.ListNavigationIndexPosition.CommitFilesComponent)
	}

	// Navigate up
	upMsg := tea.KeyPressMsg{Code: 'k'}
	m, _ = handleNonTypingUpkKeyBindingInteraction(upMsg, m)
	if m.CurrentRepoCommitFilesList.Index() != 0 {
		t.Errorf("expected index 0 after up navigation, got %d", m.CurrentRepoCommitFilesList.Index())
	}
	if m.ListNavigationIndexPosition.CommitFilesComponent != 0 {
		t.Errorf("expected ListNavigationIndexPosition 0, got %d", m.ListNavigationIndexPosition.CommitFilesComponent)
	}
}

func TestCommitFilesCheckoutAndDiscard(t *testing.T) {
	m, _, hash := setupTestModel(t)
	m.CurrentSelectedComponent = constant.CommitFilesComponentPanel
	m.CurrentDrillDownCommitHash = hash
	m.CurrentRepoCommitFilesList.Select(0)

	// Pressing 'c' should not open CommitPopUp
	m, _ = handleNonTypingcKeyBindingInteraction(m)
	if m.PopUpType == constant.CommitPopUp {
		t.Errorf("expected 'c' on CommitFilesComponentPanel not to open CommitPopUp")
	}

	// Pressing 'd' should not open GitDiscardConfirmPromptPopUp
	m, _ = handleNonTypingdKeyBindingInteraction(m)
	if m.PopUpType == constant.GitDiscardConfirmPromptPopUp {
		t.Errorf("expected 'd' on CommitFilesComponentPanel not to open GitDiscardConfirmPromptPopUp")
	}
}

func TestCommitFilesCopy(t *testing.T) {
	m, _, _ := setupTestModel(t)
	m.CurrentSelectedComponent = constant.CommitFilesComponentPanel
	m.CurrentRepoCommitFilesList.Select(0)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() {
		t.Fatalf("expected ShowPopUp to be true after pressing 'y'")
	}
	if m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Fatalf("expected PopUpType to be ChooseCopyValuePopUp, got %s", m.PopUpType)
	}
	copyModel, ok := m.PopUpModel.(*copypopup.Model)
	if !ok {
		t.Fatalf("expected PopUpModel to be *copypopup.Model")
	}
	if len(copyModel.Options.Items()) != 3 {
		t.Errorf("expected 3 copy options (relative, absolute, filename), got %d", len(copyModel.Options.Items()))
	}
}

func TestCommitFilesFilterMode(t *testing.T) {
	m, _, _ := setupTestModel(t)
	m.CurrentSelectedComponent = constant.CommitFilesComponentPanel

	m, _ = handleNonTypingFKeyBindingInteraction(m)
	if !m.IsPanelFiltering.Load() {
		t.Errorf("expected IsPanelFiltering to be true after pressing 'F'")
	}
}
