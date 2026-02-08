package tag

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// init the popup model for create tag input for ingo
func InitCreateTagPopUpModel(m *types.GittiModel, commitHash string, commitMessage string) {
	tagName := ""
	tagMessage := ""
	tagNamePlaceholder := i18n.LANGUAGEMAPPING.CreateTagPopUpNameInputPlaceHolder
	tagMessagePlaceholder := i18n.LANGUAGEMAPPING.CreateTagPopUpMessageInputPlaceHolder

	TagNameInput := textinput.New()
	TagNameInput.SetValue(tagName)
	TagNameInput.Placeholder = tagNamePlaceholder
	TagNameInput.Focus()
	TagNameInput.SetVirtualCursor(true)

	TagMessageTextAreaInput := textarea.New()
	TagMessageTextAreaInput.SetValue(tagMessage)
	TagMessageTextAreaInput.ShowLineNumbers = false
	TagMessageTextAreaInput.Placeholder = tagMessagePlaceholder
	TagMessageTextAreaInput.SetHeight(constant.TextAreaInputHeight)
	TagMessageTextAreaInput.Blur()

	popUpModel := &CreateTagPopUpModel{
		TagNameInput:            TagNameInput,
		TagMessageTextAreaInput: TagMessageTextAreaInput,
		CommitHash:              commitHash,
		CommitMessage:           commitMessage,
		CurrentActiveInputIndex: 1,
		TotalInputCount:         2,
	}
	m.PopUpModel = popUpModel
}

func InitCreateTagConfirmationPopUpModel(m *types.GittiModel, tagName string, tagMessage string, commitHash string, commitMessage string) {
	popUpModel := &CreateTagConfirmationPopUpModel{
		TagName:       tagName,
		TagMessage:    tagMessage,
		CommitHash:    commitHash,
		CommitMessage: commitMessage,
	}
	m.PopUpModel = popUpModel
}
