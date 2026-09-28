package interaction

import (
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/api"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/executor"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/component/commitfiles"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/initialize"
)

func TestClickOnCommitFilesComponentPanel(t *testing.T) {
	settings.GITTICONFIGSETTINGS = &settings.GittiDefaultConfigSettings
	i18n.InitGittiLanguageMapping("en")
	executor.InitCmdExecutor(".")
	tuiUpdateChan := make(chan interface{}, 10)
	daemonUpdateChan := make(chan string, 10)
	gitUpdateChan := make(chan string, 10)
	logger := logging.InitGittiLogging(100, make(chan string, 10), 10)
	lock := git.InitGitProcessLock(logger)
	gitCommitLog := git.InitGitCommitLog(make(chan string, 10), lock, 100, logger)
	gitOps := &api.GitOperations{
		GitCommitLog: gitCommitLog,
	}

	m := initialize.InitGittiModel(tuiUpdateChan, ".", "gitti", gitOps, logger, daemonUpdateChan, gitUpdateChan)

	// Set up layout dimensions
	m.WindowLeftPanelWidth = 50
	m.LocalBranchesComponentPanelHeight = 5
	m.ModifiedFilesComponentPanelHeight = 5
	m.CommitLogComponentPanelHeight = 10
	m.StashComponentPanelHeight = 5
	m.DetailComponentPanelHeight = 20

	// Set up commit files list
	items := []list.Item{
		commitfiles.GitCommitFileItem{CommitHash: "c1", Status: "M", FilePathname: "file1.txt"},
		commitfiles.GitCommitFileItem{CommitHash: "c1", Status: "A", FilePathname: "file2.txt"},
	}
	m.CurrentRepoCommitFilesList = list.New(items, commitfiles.GitCommitFileItemDelegate{}, 48, 10)
	m.CurrentSelectedComponent = constant.CommitFilesComponentPanel
	m.CurrentDrillDownCommitHash = "c1"

	// 1. Click on second item of commit files panel
	// Panel 0: 3 rows (top 0 to 2)
	// Panel 1: 5+2 = 7 rows (top 3 to 9)
	// Panel 2: 5+2 = 7 rows (top 10 to 16)
	// Panel 3: top = 17. Border is 17, title is 18, row 0 is 19, row 1 is 20.
	clickMsg := tea.MouseClickMsg(tea.Mouse{
		X:      10,
		Y:      20,
		Button: tea.MouseLeft,
	})

	m, _ = handleLeftMouseClick(clickMsg, m)

	if m.CurrentSelectedComponent != constant.CommitFilesComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to remain CommitFilesComponentPanel, got %s", m.CurrentSelectedComponent)
	}
	if m.CurrentRepoCommitFilesList.Index() != 1 {
		t.Fatalf("expected commit files list index to be 1, got %d", m.CurrentRepoCommitFilesList.Index())
	}

	// 2. Click on detail panel (right side, x > 51, y < 20)
	clickDetailMsg := tea.MouseClickMsg(tea.Mouse{
		X:      60,
		Y:      10,
		Button: tea.MouseLeft,
	})

	m, _ = handleLeftMouseClick(clickDetailMsg, m)

	if m.CurrentSelectedComponent != constant.DetailComponentPanel {
		t.Fatalf("expected CurrentSelectedComponent to be DetailComponentPanel, got %s", m.CurrentSelectedComponent)
	}
	if m.DetailPanelParentComponent != constant.CommitFilesComponentPanel {
		t.Fatalf("expected DetailPanelParentComponent to be CommitFilesComponentPanel, got %s", m.DetailPanelParentComponent)
	}
}
