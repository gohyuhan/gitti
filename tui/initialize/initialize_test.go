package initialize_test

import (
	"testing"

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

	// Test ReinitGittiModel
	initialize.ReinitGittiModel(m, ".", "gitti", gitOps)
	layout.LeftPanelDynamicResize(m)
}
