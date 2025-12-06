package stash

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

func InitGitStashMessagePopUpModel(m *types.GittiModel, filePathName string, stashType string) {
	stashMessageTextInput := textinput.New()
	stashMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.GitStashMessagePlaceholder
	stashMessageTextInput.Focus()
	stashMessageTextInput.SetVirtualCursor(true)
	stashMessageTextInput.SetWidth(min(constant.MaxGitStashMessagePopUpWidth, int(float64(m.Width)*0.8)) - 4)

	popUpModel := &GitStashMessagePopUpModel{
		StashMessageInput: stashMessageTextInput,
		FilePathName:      filePathName,
		StashType:         stashType,
	}

	m.PopUpModel = popUpModel
}

// for git stash output popup
func InitGitStashOperationOutputPopUpModel(m *types.GittiModel, stashOperationType string) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitStashOperationOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxGitStashOperationOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitStashOperationOutputPopUpModel{
		StashOperationType:              stashOperationType,
		GitStashOperationOutputViewport: vp,
		Spinner:                         s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	m.PopUpModel = popUpModel
}

// for git stash operation confirm prompt
func InitGitStashConfirmPromptPopUpModel(m *types.GittiModel, stashOperationType string, filePathName string, stashId string, stashMessage string) {
	popUpModel := &GitStashConfirmPromptPopUpModel{
		StashOperationType: stashOperationType,
		FilePathName:       filePathName,
		StashId:            stashId,
		StashMessage:       stashMessage,
	}
	m.PopUpModel = popUpModel
}
