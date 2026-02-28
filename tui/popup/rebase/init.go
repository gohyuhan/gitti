package rebase

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func InitGitRebaseBranchInputPopUpModel(m *types.GittiModel, remoteName string) {
	rebaseBranchNameInput := textinput.New()
	rebaseBranchNameInput.Placeholder = i18n.LANGUAGEMAPPING.RebaseBranchNameInputPlaceholder
	rebaseBranchNameInput.Focus()
	rebaseBranchNameInput.SetVirtualCursor(true)

	rebaseBranchNameInput.SetWidth(min(constant.MaxGitRebaseBranchInputPopUpWidth, int(float64(m.Width)*0.8)) - 6)
	m.PopUpModel = &GitRebaseBranchInputPopUpModel{
		BranchNameInput: rebaseBranchNameInput,
		Remote:          remoteName,
	}
}

func InitGitRebaseOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitRebaseOutputViewportHeight)
	vp.SetWidth(min(constant.MaxGitRebaseOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitRebaseOutputPopUpModel{
		GitRebaseOutputViewport: vp,
		Spinner:                 s,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.IsCancelled.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)

	m.PopUpModel = popUpModel
}
