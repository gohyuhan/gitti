package tag

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
)

type CreateTagPopUpModel struct {
	TagNameInput            textinput.Model
	TagMessageTextAreaInput textarea.Model
	CommitHash              string
	CommitMessage           string
	CurrentActiveInputIndex int
	TotalInputCount         int
}

type CreateTagConfirmationPopUpModel struct {
	TagName       string
	TagMessage    string
	CommitHash    string
	CommitMessage string
}
