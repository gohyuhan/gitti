package commit

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
)

// init the popup model for git commit
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
	CommitDescriptionTextAreaInput.SetHeight(4)
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

// init the popup model for git amend commit
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
	CommitDescriptionTextAreaInput.SetHeight(4)
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
