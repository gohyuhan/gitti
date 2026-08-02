package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/logging"
	"github.com/gohyuhan/gitti/settings"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ------------------------------------
//
//	Handle 'C' key interaction.
//	Responsibility: Handles continuation operations. Primarily triggers "git commit --continue",
//	"git rebase --continue", or "git cherry-pick --continue". Also handles the specific UI
//	suspension needed for GPG signing processes during continuations.
//
// ------------------------------------
func handleNonTypingCKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		if m.GitCommitRequireSigning && !settings.GITTICONFIGSETTINGS.OverrideSigningUISuspend {
			gitArgs := m.GitOperations.GitStateUniversalUtils.GitUniversalContinueWithSigning()
			if len(gitArgs) < 1 {
				return m, nil
			}
			return utils.SuspendGittiUIForGitOperationRequireSigning(m, gitArgs, logging.CONTINUE_COMMIT_WITH_SIGNING_OPS)
		} else {
			services.GitStateUniversalUtilsContinueService(m)
			return m, nil
		}
	}
	return m, nil
}
