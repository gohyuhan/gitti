package tag

import (
	"fmt"

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

// ------------------------------------
//
//	For Creating Tag Model
//
// ------------------------------------
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

// ------------------------------------
//
//	For Confirming Tag Creation Model
//
// ------------------------------------
func InitCreateTagConfirmationPopUpModel(m *types.GittiModel, tagName string, tagMessage string, commitHash string, commitMessage string) {
	popUpModel := &CreateTagConfirmationPopUpModel{
		TagName:       tagName,
		TagMessage:    tagMessage,
		CommitHash:    commitHash,
		CommitMessage: commitMessage,
	}
	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	For Choosing Tag Deletion Option PopUp
//
// ------------------------------------
func InitChooseDeleteTagOptionPopUpModel(m *types.GittiModel, tagName string) {
	items := make([]list.Item, 0, 2)
	items = append(items, DeleteTagOptionItem{
		Name:          fmt.Sprintf(i18n.LANGUAGEMAPPING.DeleteTagPopUpDeleteLocalTagOption, tagName),
		Info:          fmt.Sprintf(i18n.LANGUAGEMAPPING.DeleteTagPopUpDeleteLocalTagOptionInfo, tagName),
		DeleteTagType: git.TAGDELETELOCAL,
	})
	items = append(items, DeleteTagOptionItem{
		Name:          fmt.Sprintf(i18n.LANGUAGEMAPPING.DeleteTagPopUpDeleteRemoteTagOption, tagName),
		Info:          fmt.Sprintf(i18n.LANGUAGEMAPPING.DeleteTagPopUpDeleteRemoteTagOptionInfo, tagName),
		DeleteTagType: git.TAGDELETEREMOTE,
	})
	width := (min(constant.MaxChooseDeleteTagOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	dL := list.New(items, DeleteTagOptionDelegate{}, width, constant.PopUpChooseDeleteTagOptionHeight)
	dL.SetShowPagination(false)
	dL.SetShowStatusBar(false)
	dL.SetFilteringEnabled(false)
	dL.SetShowTitle(false)

	// Custom Help Model for Count Display
	dL.SetShowHelp(true)
	dL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	dL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	dL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &dL, constant.MaxChooseDeleteTagOptionPopUpWidth)

	m.PopUpModel = &ChooseDeleteTagOptionPopUpModel{
		TagName:          tagName,
		DeleteOptionList: dL,
	}
}

// ------------------------------------
//
//	For Choosing Remote for Remote Tag Deletion PopUp
//
// ------------------------------------
func InitChooseRemoteForDeleteRemoteTagPopUpModel(m *types.GittiModel, remoteList []git.GitRemoteInfo, tagName string, deleteOptionType string) {
	items := make([]list.Item, 0, len(remoteList))
	for _, remote := range remoteList {
		items = append(items, GitRemoteForDeleteRemoteTagItem(remote))
	}
	width := (min(constant.MaxChooseRemoteForDeleteRemoteTagPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rL := list.New(items, GitRemoteForDeleteRemoteTagItemDelegate{}, width, constant.PopUpChooseRemoteHeight)
	rL.SetShowPagination(false)
	rL.SetShowStatusBar(false)
	rL.SetFilteringEnabled(false)
	rL.SetShowTitle(false)

	// Custom Help Model for Count Display
	rL.SetShowHelp(true)
	rL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	rL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	rL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &rL, constant.MaxChooseRemoteForDeleteRemoteTagPopUpWidth)

	m.PopUpModel = &ChooseRemoteForDeleteRemoteTagPopUpModel{
		RemoteList:       rL,
		TagName:          tagName,
		DeleteOptionType: deleteOptionType,
	}
}

// ------------------------------------
//
//	For Tag Deletion Output PopUp
//
// ------------------------------------
func InitDeleteTagOutputPopUpModel(m *types.GittiModel, tagName string) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpDeleteTagOutputViewportHeight)
	vp.SetWidth(min(constant.MaxDeleteTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &DeleteTagOutputPopUpModel{
		TagName:                 tagName,
		DeleteTagOutputViewport: vp,
		Spinner:                 s,
		CancelFunc:              nil,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)

	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	For Tag Push Option PopUp
//
// ------------------------------------
func InitChoosePushTagOptionPopUpModel(m *types.GittiModel, remoteName string, tagName string) {
	items := make([]list.Item, 0, 4)
	items = append(items, PushTagOptionItem{
		Name:        fmt.Sprintf(i18n.LANGUAGEMAPPING.PushTagPopUpPushTagOption, tagName),
		Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.PushTagPopUpPushTagOptionInfo, tagName),
		PushTagType: git.TAGPUSH,
	})
	items = append(items, PushTagOptionItem{
		Name:        i18n.LANGUAGEMAPPING.PushTagPopUpPushAllTagOption,
		Info:        i18n.LANGUAGEMAPPING.PushTagPopUpPushAllTagOptionInfo,
		PushTagType: git.TAGPUSHALL,
	})
	items = append(items, PushTagOptionItem{
		Name:        fmt.Sprintf(i18n.LANGUAGEMAPPING.PushTagPopUpPushForceTagOption, tagName),
		Info:        fmt.Sprintf(i18n.LANGUAGEMAPPING.PushTagPopUpPushForceTagOptionInfo, tagName),
		PushTagType: git.TAGPUSHFORCE,
	})
	items = append(items, PushTagOptionItem{
		Name:        i18n.LANGUAGEMAPPING.PushTagPopUpPushAllForceTagOption,
		Info:        i18n.LANGUAGEMAPPING.PushTagPopUpPushAllForceTagOptionInfo,
		PushTagType: git.TAGPUSHALLFORCE,
	})
	width := (min(constant.MaxChoosePushTagOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	cPTOL := list.New(items, PushTagOptionDelegate{}, width, constant.PopUpChoosePushTagOptionHeight)
	cPTOL.SetShowPagination(false)
	cPTOL.SetShowStatusBar(false)
	cPTOL.SetFilteringEnabled(false)
	cPTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cPTOL.SetShowHelp(true)
	cPTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cPTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cPTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cPTOL, constant.MaxChoosePushTagOptionPopUpWidth)

	m.PopUpModel = &ChoosePushTagOptionPopUpModel{
		TagName:        tagName,
		RemoteName:     remoteName,
		PushOptionList: cPTOL,
	}
}

// ------------------------------------
//
//	For Tag Push Output PopUp
//
// ------------------------------------
func InitPushTagOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpPushTagOutputViewportHeight)
	vp.SetWidth(min(constant.MaxPushTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &PushTagOutputPopUpModel{
		PushTagOutputViewport: vp,
		Spinner:               s,
		CancelFunc:            nil,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)

	m.PopUpModel = popUpModel
}
