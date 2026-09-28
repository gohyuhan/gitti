package nontyping

import (
	"testing"

	"charm.land/bubbles/v2/list"
	"github.com/atotto/clipboard"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/component/branch"
	"github.com/gohyuhan/gitti/tui/component/commitlog"
	"github.com/gohyuhan/gitti/tui/component/files"
	"github.com/gohyuhan/gitti/tui/component/reflog"
	"github.com/gohyuhan/gitti/tui/component/remote"
	"github.com/gohyuhan/gitti/tui/component/stash"
	"github.com/gohyuhan/gitti/tui/component/tag"
	"github.com/gohyuhan/gitti/tui/component/worktree"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

func newTestModel() *types.GittiModel {
	i18n.InitGittiLanguageMapping("EN")
	logChan := make(chan string, 10)
	logger := logging.InitGittiLogging(50, logChan, 5)

	m := &types.GittiModel{
		GittiLogger: logger,
		Width:       100,
		Height:      50,
	}

	return m
}

func assertClipboard(t *testing.T, expected string) {
	t.Helper()
	got, err := clipboard.ReadAll()
	if err != nil {
		t.Logf("Skipping clipboard assertion in headless environment: %v", err)
		return
	}
	if got != expected {
		t.Errorf("expected '%s' in clipboard, got '%s'", expected, got)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_EmptyListNoPanic(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.CommitLogOrRefLogComponentPanel
	m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_COMMITLOG

	// Should not panic on empty list
	m, cmd := handleNonTypingyKeyBindingInteraction(m)
	if cmd != nil {
		t.Errorf("expected nil cmd, got %v", cmd)
	}
	if m.ShowPopUp.Load() {
		t.Errorf("expected popup not to show on empty list")
	}
}

func TestHandleNonTypingyKeyBindingInteraction_CommitLog(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.CommitLogOrRefLogComponentPanel
	m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_COMMITLOG

	items := []list.Item{
		commitlog.GitCommitLogItem{
			Hash:    "abcdef1234567890abcdef1234567890abcdef12",
			Message: "feat: initial commit",
			Author:  "Test User",
		},
	}
	m.CurrentRepoCommitLogInfoList = list.New(items, commitlog.GitCommitLogItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_Branch(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel
	m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_LOCAL_BRANCH

	items := []list.Item{
		branch.GitBranchItem{
			BranchName:   "feature/test-branch",
			IsCheckedOut: true,
		},
	}
	m.CurrentRepoBranchesInfoList = list.New(items, branch.GitBranchItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_ModifiedFiles(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.ModifiedFilesComponentPanel

	items := []list.Item{
		files.GitModifiedFilesItem{
			FilePathname:    "src/main.go",
			NewFilePathname: "src/main.go",
			IndexState:      "M",
			WorkTree:        " ",
		},
	}
	m.CurrentRepoModifiedFilesInfoList = list.New(items, files.GitModifiedFilesItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_Tag(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel
	m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_TAG

	items := []list.Item{
		tag.GitTagItem{
			TagName: "v1.0.0",
		},
	}
	m.CurrentRepoTagInfoList = list.New(items, tag.GitTagItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_Remote(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel
	m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_REMOTE

	items := []list.Item{
		remote.GitRemoteItem{
			Name: "origin",
			Url:  "git@github.com:org/repo.git",
		},
	}
	m.CurrentRepoRemoteInfoList = list.New(items, remote.GitRemoteItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_Worktree(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.LocalBranchOrTagOrRemoteOrWorktreeComponentPanel
	m.CurrentLocalBranchOrTagOrRemoteOrWorktreeComponentShowing = constant.SHOW_WORKTREE

	items := []list.Item{
		worktree.GitWorktreeItem{
			WorktreePath: "/tmp/my-worktree",
		},
	}
	m.CurrentRepoWorktreeInfoList = list.New(items, worktree.GitWorktreeItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_Stash(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.StashComponentPanel

	items := []list.Item{
		stash.GitStashItem{
			Id:      "stash@{0}",
			Message: "WIP on main",
		},
	}
	m.CurrentRepoStashInfoList = list.New(items, stash.GitStashItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_RefLog(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.CommitLogOrRefLogComponentPanel
	m.CurrentCommitLogOrRefLogComponentShowing = constant.SHOW_REFLOG

	items := []list.Item{
		reflog.GitRefLogItem{
			Hash:     "1234567890abcdef1234567890abcdef12345678",
			InfoDesc: "checkout: moving from main to feat",
		},
	}
	m.CurrentRepoRefLogInfoList = list.New(items, reflog.GitRefLogItemDelegate{}, 100, 20)

	m, _ = handleNonTypingyKeyBindingInteraction(m)
	if !m.ShowPopUp.Load() || m.PopUpType != constant.ChooseCopyValuePopUp {
		t.Errorf("expected ChooseCopyValuePopUp to show, got show=%v type=%s", m.ShowPopUp.Load(), m.PopUpType)
	}
}

func TestHandleNonTypingyKeyBindingInteraction_LineEditing(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.DetailComponentPanel
	m.IsLineEditingState.Store(true)
	m.DetailPanelViewportOGStringArray = []string{
		"@@ -1,3 +1,4 @@",
		"-oldLine",
		"+newLine",
	}
	m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = 2

	m, cmd := handleNonTypingyKeyBindingInteraction(m)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from StartCopy")
	}
	msg := cmd()
	finishedMsg, ok := msg.(types.CopyFinishedMsg)
	if !ok {
		t.Fatalf("expected types.CopyFinishedMsg, got %T", msg)
	}
	if finishedMsg.Value != "+newLine" {
		t.Errorf("expected copied value '+newLine', got '%s'", finishedMsg.Value)
	}
	assertClipboard(t, "+newLine")
}

func TestHandleNonTypingyKeyBindingInteraction_LineEditing_EmptyLine(t *testing.T) {
	m := newTestModel()
	m.CurrentSelectedComponent = constant.DetailComponentPanel
	m.IsLineEditingState.Store(true)
	m.DetailPanelViewportOGStringArray = []string{
		"@@ -1,3 +1,4 @@",
		"",
	}
	m.LineEditingIndexPositionAndInfo.DetailPanelViewportActualCurrentIndex = 1

	// Seed with non-empty string first
	_ = clipboard.WriteAll("some-previous-content")

	m, cmd := handleNonTypingyKeyBindingInteraction(m)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from StartCopy")
	}
	msg := cmd()
	finishedMsg, ok := msg.(types.CopyFinishedMsg)
	if !ok {
		t.Fatalf("expected types.CopyFinishedMsg, got %T", msg)
	}
	if finishedMsg.Value != "" {
		t.Errorf("expected copied value '', got '%s'", finishedMsg.Value)
	}
	assertClipboard(t, "")
}
