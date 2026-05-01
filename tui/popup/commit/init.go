package commit

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// ----------------------------------
//
//	init the popup model for git commit
//
// ----------------------------------
func InitGitCommitPopUpModel(m *types.GittiModel) {
	commitMsg := ""
	commitDesc := ""
	commitMsgPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolder
	commitDescPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolder

	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.SetValue(commitMsg)
	CommitMessageTextInput.Placeholder = commitMsgPlaceholder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.SetVirtualCursor(true)

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.SetValue(commitDesc)
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = commitDescPlaceholder
	CommitDescriptionTextAreaInput.DynamicHeight = true
	CommitDescriptionTextAreaInput.MinHeight = constant.TextAreaInputMinHeight
	CommitDescriptionTextAreaInput.MaxHeight = constant.TextAreaInputMaxHeight
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitCommitOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitCommitPopUpModel{
		IsAmendCommit:            false,
		MessageTextInput:         CommitMessageTextInput,
		DescriptionTextAreaInput: CommitDescriptionTextAreaInput,
		TotalInputCount:          2,
		CurrentActiveInputIndex:  1,
		GitCommitOutputViewport:  vp,
		Spinner:                  s,
	}
	popUpModel.InitialCommitStarted.Store(false)
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// ----------------------------------
//
//	init the popup model for git amend commit
//
// ----------------------------------
func InitGitAmendCommitPopUpModel(m *types.GittiModel) {
	commitMsgAndDesc := m.GitOperations.GitCommit.GetLatestCommitMsgAndDesc()
	commitMsg := commitMsgAndDesc.Message
	commitDesc := commitMsgAndDesc.Description
	commitMsgPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpMessageInputPlaceHolderAmendVersion
	commitDescPlaceholder := i18n.LANGUAGEMAPPING.CommitPopUpCommitDescriptionInputPlaceHolderAmendVersion

	CommitMessageTextInput := textinput.New()
	CommitMessageTextInput.SetValue(commitMsg)
	CommitMessageTextInput.Placeholder = commitMsgPlaceholder
	CommitMessageTextInput.Focus()
	CommitMessageTextInput.SetVirtualCursor(true)

	CommitDescriptionTextAreaInput := textarea.New()
	CommitDescriptionTextAreaInput.SetValue(commitDesc)
	CommitDescriptionTextAreaInput.ShowLineNumbers = false
	CommitDescriptionTextAreaInput.Placeholder = commitDescPlaceholder
	CommitDescriptionTextAreaInput.DynamicHeight = true
	CommitDescriptionTextAreaInput.MinHeight = constant.TextAreaInputMinHeight
	CommitDescriptionTextAreaInput.MaxHeight = constant.TextAreaInputMaxHeight
	CommitDescriptionTextAreaInput.MoveToEnd()
	CommitDescriptionTextAreaInput.Blur()

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitAmendCommitOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxAmendCommitPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitAmendCommitPopUpModel{
		IsAmendCommit:                true,
		MessageTextInput:             CommitMessageTextInput,
		DescriptionTextAreaInput:     CommitDescriptionTextAreaInput,
		TotalInputCount:              2,
		CurrentActiveInputIndex:      1,
		GitAmendCommitOutputViewport: vp,
		Spinner:                      s,
	}
	popUpModel.InitialCommitStarted.Store(false)
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// ----------------------------------
//
//	for reset latest commit option list popup
//
// ----------------------------------
func InitGitResetLatestCommitTypeOptionPopUpModel(m *types.GittiModel) {
	gitResetLatestCommitTypeOption := []GitResetLatestCommitTypeOptionItem{
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetSoft,
			Info:      i18n.LANGUAGEMAPPING.GitResetSoftInfo,
			ResetType: git.RESETSOFT,
		},
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetHard,
			Info:      i18n.LANGUAGEMAPPING.GitResetHardInfo,
			ResetType: git.RESETHARD,
		},
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetMixed,
			Info:      i18n.LANGUAGEMAPPING.GitResetMixedInfo,
			ResetType: git.RESETMIXED,
		},
	}

	items := make([]list.Item, 0, len(gitResetLatestCommitTypeOption))
	for _, resetOption := range gitResetLatestCommitTypeOption {
		items = append(items, GitResetLatestCommitTypeOptionItem(resetOption))
	}

	width := (min(constant.MaxGitResetLatestCommitTypeOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gRLCTOL := list.New(items, GitResetLatestCommitTypeOptionDelegate{}, width, constant.PopUpGitResetLatestCommitTypeOptionHeight)
	gRLCTOL.SetShowPagination(false)
	gRLCTOL.SetShowStatusBar(false)
	gRLCTOL.SetFilteringEnabled(false)
	gRLCTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gRLCTOL.SetShowHelp(true)
	gRLCTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gRLCTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gRLCTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gRLCTOL, constant.MaxGitResetLatestCommitTypeOptionPopUpWidth)

	popUpModel := &GitResetLatestCommitTypeOptionPopUpModel{
		ResetLatestCommitTypeOptionList: gRLCTOL,
	}

	m.PopUpModel = popUpModel
}

// ----------------------------------
//
//	for git reset latest commit confirmation prompt
//
// ----------------------------------
func InitGitResetLatestCommitConfirmPromptPopUpModel(m *types.GittiModel, resetType string) {
	popUpModel := &GitResetLatestCommitConfirmPromptPopUpModel{
		GitResetLatestCommitType: resetType,
	}
	m.PopUpModel = popUpModel
}

// ----------------------------------
//
//	for reset selected commit option list popup
//
// ----------------------------------
func InitGitResetToSelectedCommitTypeOptionPopUpModel(m *types.GittiModel, selectedCommitHash string, commitInfoMessage string, commitInfoAuthor string) {
	gitResetToSelectedCommitTypeOption := []GitResetToSelectedCommitTypeOptionItem{
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetSoft,
			Info:      i18n.LANGUAGEMAPPING.GitResetSoftInfo,
			ResetType: git.RESETSOFT,
		},
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetHard,
			Info:      i18n.LANGUAGEMAPPING.GitResetHardInfo,
			ResetType: git.RESETHARD,
		},
		{
			Name:      i18n.LANGUAGEMAPPING.GitResetMixed,
			Info:      i18n.LANGUAGEMAPPING.GitResetMixedInfo,
			ResetType: git.RESETMIXED,
		},
	}

	items := make([]list.Item, 0, len(gitResetToSelectedCommitTypeOption))
	for _, resetOption := range gitResetToSelectedCommitTypeOption {
		items = append(items, GitResetToSelectedCommitTypeOptionItem(resetOption))
	}

	width := (min(constant.MaxGitResetToSelectedCommitTypeOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	gRSCTOL := list.New(items, GitResetToSelectedCommitTypeOptionDelegate{}, width, constant.PopUpGitResetToSelectedCommitTypeOptionHeight)
	gRSCTOL.SetShowPagination(false)
	gRSCTOL.SetShowStatusBar(false)
	gRSCTOL.SetFilteringEnabled(false)
	gRSCTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	gRSCTOL.SetShowHelp(true)
	gRSCTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	gRSCTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	gRSCTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &gRSCTOL, constant.MaxGitResetToSelectedCommitTypeOptionPopUpWidth)

	popUpModel := &GitResetToSelectedCommitTypeOptionPopUpModel{
		ResetToSelectedCommitTypeOptionList: gRSCTOL,
		SelectedCommitHash:                  selectedCommitHash,
		CommitInfoMessage:                   commitInfoMessage,
		CommitInfoAuthor:                    commitInfoAuthor,
	}

	m.PopUpModel = popUpModel
}

// ----------------------------------
//
//	for git reset selected commit confirmation prompt
//
// ----------------------------------
func InitGitResetToSelectedCommitConfirmPromptPopUpModel(m *types.GittiModel, resetType string, selectedCommitHash string, commitInfoMessage string, commitInfoAuthor string) {
	popUpModel := &GitResetToSelectedCommitConfirmPromptPopUpModel{
		GitResetToSelectedCommitType: resetType,
		SelectedCommitHash:           selectedCommitHash,
		CommitInfoMessage:            commitInfoMessage,
		CommitInfoAuthor:             commitInfoAuthor,
	}
	m.PopUpModel = popUpModel
}
