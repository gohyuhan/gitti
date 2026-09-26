package services

import (
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/utils"
)

// ------------------------------------
//
//	Open the web page of a commit in the browser. The git lookups run off the UI thread,
//	because checking that the commit is on the remote can be slow in a large repo.
//
// ------------------------------------
func OpenCommitWebPageService(m *types.GittiModel, commitHash string) {
	checkedOutBranchName := m.CheckOutBranch
	go func() {
		webPageURL, err := git.GetCommitWebPageURL(commitHash, checkedOutBranchName)
		openWebPage(m, webPageURL, err)
	}()
}

// ------------------------------------
//
//	Open the branch page (git.WEBPAGEBRANCH) or the new pull request page
//	(git.WEBPAGEPULLREQUEST) of a local branch in the browser
//
// ------------------------------------
func OpenBranchWebPageService(m *types.GittiModel, branchName string, pageType string) {
	go func() {
		webPageURL, err := git.GetBranchWebPageURL(branchName, pageType)
		openWebPage(m, webPageURL, err)
	}()
}

// ------------------------------------
//
//	Log why the web page cannot open, or log its URL and open it
//
// ------------------------------------
func openWebPage(m *types.GittiModel, webPageURL string, err error) {
	if err != nil {
		// the host in the message comes from the remote URL, so keep terminal escape sequences out of the log panel
		m.GittiLogger.RegisterNewLog(logging.OPEN_IN_BROWSER_OPS, "", logging.WARN, utils.EscapeControlCharacters(err.Error()), false)
		return
	}
	m.GittiLogger.RegisterNewLog(logging.OPEN_IN_BROWSER_OPS, webPageURL, logging.INFO, "", false)
	utils.OpenBrowser(webPageURL)
}
