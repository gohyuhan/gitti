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
	newWorktreeNameInput := textinput.New()
	newWorktreeNameInput.Placeholder = i18n.LANGUAGEMAPPING.CreateNewWorktreePrompt
	newWorktreeNameInput.Focus()
	newWorktreeNameInput.SetVirtualCursor(true)

	newWorktreeBranchInput := textinput.New()
	newWorktreeBranchInput.Placeholder = i18n.LANGUAGEMAPPING.NewWorktreeBranchPrompt
	newWorktreeBranchInput.Blur()
	newWorktreeBranchInput.SetVirtualCursor(true)

	newWorktreeNameInput.SetWidth(min(constant.MaxWorktreeAddNewWorktreePopUpWidth, int(float64(m.Width)*0.8)) - 6)
	newWorktreeBranchInput.SetWidth(min(constant.MaxWorktreeAddNewWorktreePopUpWidth, int(float64(m.Width)*0.8)) - 6)
	m.PopUpModel = &WorktreeAddNewWorktreePopUpModel{
		WorktreeNameTextInput:       newWorktreeNameInput,
		WorktreeBranchNameTextInput: newWorktreeBranchInput,
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

// ------------------------------------
//
//	Initialize the lock-reason input popup model with a single focused text input
//	sized to fit the terminal width, tied to the worktree at worktreePath.
//
// ------------------------------------
func InitWorktreeLockReasonInputPopUpModel(m *types.GittiModel, worktreePath string) {
	worktreeLockReasonInput := textinput.New()
	worktreeLockReasonInput.Placeholder = i18n.LANGUAGEMAPPING.WorktreeLockReasonPrompt
	worktreeLockReasonInput.Focus()
	worktreeLockReasonInput.SetVirtualCursor(true)

	worktreeLockReasonInput.SetWidth(min(constant.MaxWorktreeLockReasonInputPopUpWidth, int(float64(m.Width)*0.8)) - 6)
	m.PopUpModel = &WorktreeLockReasonInputPopUpModel{
		WorktreeLockReasonTextInput: worktreeLockReasonInput,
		WorktreePath:                worktreePath,
	}
}

// ------------------------------------
//
//	Initialize the remove-worktree confirmation popup model tied to the worktree
//	at worktreePath.
//
// ------------------------------------
func InitWorktreeRemoveWorktreeConfirmationPopUpModel(m *types.GittiModel, worktreePath string) {
	m.PopUpModel = &WorktreeRemoveWorktreeConfirmationPopUpModel{
		WorktreePath: worktreePath,
	}
}
