package worktree

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
//	Initialize the add-new-worktree input popup model with two focused-aware text
//	inputs (worktree name and optional branch name) sized to fit the terminal
//	width, defaulting focus to the first input.
//
// ------------------------------------
func InitWorktreeAddNewWorktreePopUpModel(m *types.GittiModel) {
	NewWorktreeNameInput := textinput.New()
	NewWorktreeNameInput.Placeholder = i18n.LANGUAGEMAPPING.CreateNewWorktreePrompt
	NewWorktreeNameInput.Focus()
	NewWorktreeNameInput.SetVirtualCursor(true)

	NewWorktreeBranchInput := textinput.New()
	NewWorktreeBranchInput.Placeholder = i18n.LANGUAGEMAPPING.NewWorktreeBranchPrompt
	NewWorktreeBranchInput.Blur()
	NewWorktreeBranchInput.SetVirtualCursor(true)

	NewWorktreeNameInput.SetWidth(min(constant.MaxWorktreeAddNewWorktreePopUpWidth, int(float64(m.Width)*0.8)) - 6)
	NewWorktreeBranchInput.SetWidth(min(constant.MaxWorktreeAddNewWorktreePopUpWidth, int(float64(m.Width)*0.8)) - 6)
	m.PopUpModel = &WorktreeAddNewWorktreePopUpModel{
		WorktreeNameTextInput:       NewWorktreeNameInput,
		WorktreeBranchNameTextInput: NewWorktreeBranchInput,
		TotalInputCount:             2,
		CurrentActiveInputIndex:     1,
	}
}

// ------------------------------------
//
//	Initialize the add-new-worktree output popup model with a soft-wrap viewport,
//	a dot spinner, and all atomic state flags reset to false.
//
// ------------------------------------
func InitWorktreeAddNewWorktreeOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpWorktreeAddNewWorktreeOutputViewportHeight)
	vp.SetWidth(min(constant.MaxWorktreeAddNewWorktreeOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &WorktreeAddNewWorktreeOutputPopUpModel{
		AddNewWorktreeOutputViewport: vp,
		Spinner:                      s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}
