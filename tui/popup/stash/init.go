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

// ------------------------------------
//
//	Initialize the stash message prompt popup. Creates a focused text input with
//	a localized placeholder, sized to 80% of the terminal width (minus inner
//	padding). filePathName is set when stashing a single file; stashType
//	distinguishes stash-file from stash-all operations.
//
// ------------------------------------
func InitGitStashMessagePopUpModel(m *types.GittiModel, filePathName string, stashType string) {
	stashMessageTextInput := textinput.New()
	stashMessageTextInput.Placeholder = i18n.LANGUAGEMAPPING.GitStashMessagePlaceholder
	stashMessageTextInput.Focus()
	stashMessageTextInput.SetVirtualCursor(true)
	stashMessageTextInput.SetWidth(min(constant.MaxGitStashMessagePopUpWidth, int(float64(m.Width)*0.8)) - 6)

	popUpModel := &GitStashMessagePopUpModel{
		StashMessageInput: stashMessageTextInput,
		FilePathName:      filePathName,
		StashType:         stashType,
	}

	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	Initialize the stash operation output popup. Creates a soft-wrap viewport
//	sized to a fixed height and 80% of terminal width (minus padding), and a dot
//	spinner. All atomic state flags (IsProcessing, HasError, ProcessSuccess) are
//	reset to false. stashOperationType selects the localized title and processing
//	text shown during rendering.
//
// ------------------------------------
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

// ------------------------------------
//
//	Initialize the stash confirmation prompt popup with the target operation type,
//	file path (single-file stash), stash ID, and stash message. The render
//	function uses these fields to build the localized confirmation string.
//
// ------------------------------------
func InitGitStashConfirmPromptPopUpModel(m *types.GittiModel, stashOperationType string, filePathName string, stashId string, stashMessage string) {
	popUpModel := &GitStashConfirmPromptPopUpModel{
		StashOperationType: stashOperationType,
		FilePathName:       filePathName,
		StashId:            stashId,
		StashMessage:       stashMessage,
	}
	m.PopUpModel = popUpModel
}
