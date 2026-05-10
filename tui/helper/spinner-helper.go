package helper

import (
	"github.com/gohyuhan/gitti/tui/constant"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	commitPopUp "github.com/gohyuhan/gitti/tui/popup/commit"
	interactiverebasePopUp "github.com/gohyuhan/gitti/tui/popup/interactive-rebase"
	pullPopUp "github.com/gohyuhan/gitti/tui/popup/pull"
	pushPopUp "github.com/gohyuhan/gitti/tui/popup/push"
	rebasePopUp "github.com/gohyuhan/gitti/tui/popup/rebase"
	stashPopUp "github.com/gohyuhan/gitti/tui/popup/stash"
	tagPopUp "github.com/gohyuhan/gitti/tui/popup/tag"
	"github.com/gohyuhan/gitti/tui/types"

	tea "charm.land/bubbletea/v2"
)

// ------------------------------------
//
//	Tick the spinner for every popup that is currently processing. Dispatches
//	to the correct popup model type based on m.PopUpType, calls Spinner.Update,
//	and appends the resulting Cmd to cmds.
//
// ------------------------------------
func UpdateSpinner(m *types.GittiModel, msg tea.Msg, cmds []tea.Cmd) []tea.Cmd {
	if m.ShowPopUp.Load() {
		switch m.PopUpType {
		case constant.CommitPopUp:
			if commitPopup, ok := m.PopUpModel.(*commitPopUp.GitCommitPopUpModel); ok && commitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				commitPopup.Spinner, cmd = commitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.AmendCommitPopUp:
			if amendCommitPopup, ok := m.PopUpModel.(*commitPopUp.GitAmendCommitPopUpModel); ok && amendCommitPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				amendCommitPopup.Spinner, cmd = amendCommitPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitRemotePushPopUp:
			if pushPopup, ok := m.PopUpModel.(*pushPopUp.GitRemotePushPopUpModel); ok && pushPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pushPopup.Spinner, cmd = pushPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitPullOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*pullPopUp.GitPullOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.SwitchBranchOutputPopUp:
			if pullPopup, ok := m.PopUpModel.(*branchPopUp.SwitchBranchOutputPopUpModel); ok && pullPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				pullPopup.Spinner, cmd = pullPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitStashOperationOutputPopUp:
			if stashPopup, ok := m.PopUpModel.(*stashPopUp.GitStashOperationOutputPopUpModel); ok && stashPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				stashPopup.Spinner, cmd = stashPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitDeleteBranchOutputPopUp:
			if branchPopup, ok := m.PopUpModel.(*branchPopUp.GitDeleteBranchOutputPopUpModel); ok && branchPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				branchPopup.Spinner, cmd = branchPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.CreateBranchBasedOnRemoteOutputPopUp:
			if branchPopup, ok := m.PopUpModel.(*branchPopUp.CreateBranchBasedOnRemoteOutputPopUpModel); ok && branchPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				branchPopup.Spinner, cmd = branchPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.DeleteTagOutputPopUp:
			if tagPopup, ok := m.PopUpModel.(*tagPopUp.DeleteTagOutputPopUpModel); ok && tagPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				tagPopup.Spinner, cmd = tagPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.PushTagOutputPopUp:
			tagPopUp, ok := m.PopUpModel.(*tagPopUp.PushTagOutputPopUpModel)
			if ok {
				var cmd tea.Cmd
				tagPopUp.Spinner, cmd = tagPopUp.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.FetchTagOutputPopUp:
			if tagPopup, ok := m.PopUpModel.(*tagPopUp.FetchTagOutputPopUpModel); ok && tagPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				tagPopup.Spinner, cmd = tagPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.GitRebaseOutputPopUp:
			if rebasePopup, ok := m.PopUpModel.(*rebasePopUp.GitRebaseOutputPopUpModel); ok && rebasePopup.IsProcessing.Load() {
				var cmd tea.Cmd
				rebasePopup.Spinner, cmd = rebasePopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.BranchMergeOutputPopUp:
			if branchPopup, ok := m.PopUpModel.(*branchPopUp.BranchMergeOutputPopUpModel); ok && branchPopup.IsProcessing.Load() {
				var cmd tea.Cmd
				branchPopup.Spinner, cmd = branchPopup.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		case constant.InteractiveRebaseFixupSquashOutputPopUp:
			if interactiveRebaseFixupSquash, ok := m.PopUpModel.(*interactiverebasePopUp.InteractiveRebaseFixupSquashOutputPopUpModel); ok && interactiveRebaseFixupSquash.IsProcessing.Load() {
				var cmd tea.Cmd
				interactiveRebaseFixupSquash.Spinner, cmd = interactiveRebaseFixupSquash.Spinner.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return cmds
}
