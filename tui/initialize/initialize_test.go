package initialize_test

import (
	"testing"

	"charm.land/bubbles/v2/list"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/initialize"
	"github.com/gohyuhan/gitti/tui/layout"
)

func TestInitGittiModelCommitFilesListInitialized(t *testing.T) {
	settings.GITTICONFIGSETTINGS = &settings.GittiDefaultConfigSettings
	executor.InitCmdExecutor(".")
	tuiUpdateChan := make(chan interface{}, 10)
	daemonUpdateChan := make(chan string, 10)
	gitUpdateChan := make(chan string, 10)
	logger := logging.InitGittiLogging(100, make(chan string, 10), 10)
	gitOps := &api.GitOperations{}

	m := initialize.InitGittiModel(tuiUpdateChan, ".", "gitti", gitOps, logger, daemonUpdateChan, gitUpdateChan)
	if m == nil {
		t.Fatal("expected InitGittiModel to return non-nil model")
	}

	m.WindowLeftPanelWidth = 50
	m.CommitLogComponentPanelHeight = 20

	// Must not panic on LeftPanelDynamicResize
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic during LeftPanelDynamicResize: %v", r)
		}
	}()

	layout.LeftPanelDynamicResize(m)

	// Set cache values to verify ReinitGittiModel clears them
	m.CommitFilesCache = make([]list.Item, 1)
	m.CommitFilesCommitHash = "abcdef123"

	// Test ReinitGittiModel
	initialize.ReinitGittiModel(m, ".", "gitti", gitOps)
	if m.CommitFilesCache != nil {
		t.Errorf("expected CommitFilesCache to be nil after ReinitGittiModel, got %v", m.CommitFilesCache)
	}
	if m.CommitFilesCommitHash != "" {
		t.Errorf("expected CommitFilesCommitHash to be empty after ReinitGittiModel, got %q", m.CommitFilesCommitHash)
	}
	layout.LeftPanelDynamicResize(m)
}
