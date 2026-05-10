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
//	Initialize the create tag popup. Creates a focused tag-name text input and an
//	unfocused dynamic-height message textarea, both with localized placeholders.
//	CurrentActiveInputIndex is set to 1 (name field focused first).
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
	TagMessageTextAreaInput.DynamicHeight = true
	TagMessageTextAreaInput.MinHeight = constant.TextAreaInputMinHeight
	TagMessageTextAreaInput.MaxHeight = constant.TextAreaInputMaxHeight
	TagMessageTextAreaInput.MaxContentHeight = 9999
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
//	Initialize the tag creation confirmation popup with the tag name, optional
//	message, and the target commit hash and message. The render function uses
//	these fields to build the localized confirmation string.
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
//	Initialize the delete tag scope selection popup. Builds a two-option list
//	(delete local, delete remote) with localized names and info strings. Filtering
//	and pagination are hidden; an item-count help key is attached.
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
//	Initialize the remote selection popup for remote tag deletion. Builds a list
//	from the provided remotes (each showing name and URL). Filtering and
//	pagination are hidden; an item-count help key is attached.
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
//	Initialize the delete tag output popup. Creates a soft-wrap viewport sized
//	to a fixed height and 80% of terminal width (minus padding), and a dot
//	spinner. All atomic state flags (IsProcessing, HasError, ProcessSuccess,
//	IsCancelled) are reset to false.
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
//	Initialize the push tag option popup for the given remote and tag. Builds a
//	four-option list: push tag, push all tags, force-push tag, force-push all
//	tags. Filtering and pagination are hidden; an item-count help key is attached.
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
//	Initialize the push tag output popup. Creates a soft-wrap viewport sized to
//	a fixed height and 80% of terminal width (minus padding), and a dot spinner.
//	All atomic state flags (IsProcessing, HasError, ProcessSuccess, IsCancelled)
//	are reset to false.
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

// ------------------------------------
//
//	Initialize the fetch tag option popup for the given remote. Builds a
//	four-option list: fetch tags, fetch and overwrite existing tags, fetch and
//	prune deleted remote tags, or mirror all tags. Filtering and pagination are
//	hidden; an item-count help key is attached.
//
// ------------------------------------
func InitChooseFetchTagOptionPopUpModel(m *types.GittiModel, remoteName string) {
	items := make([]list.Item, 0, 4)
	items = append(items, FetchTagOptionItem{
		Name:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchTagOption,
		Info:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchTagOptionInfo,
		FetchTagType: git.TAGFETCH,
	})
	items = append(items, FetchTagOptionItem{
		Name:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchOverwriteTagOption,
		Info:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchOverwriteTagOptionInfo,
		FetchTagType: git.TAGFETCHOVERWRITE,
	})
	items = append(items, FetchTagOptionItem{
		Name:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchPruneTagOption,
		Info:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchPruneTagOptionInfo,
		FetchTagType: git.TAGFETCHPRUNE,
	})
	items = append(items, FetchTagOptionItem{
		Name:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchMirrorTagOption,
		Info:         i18n.LANGUAGEMAPPING.FetchTagPopUpFetchMirrorTagOptionInfo,
		FetchTagType: git.TAGFETCHMIRROR,
	})
	width := (min(constant.MaxChooseFetchTagOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	cFTOL := list.New(items, FetchTagOptionDelegate{}, width, constant.PopUpChooseFetchTagOptionHeight)
	cFTOL.SetShowPagination(false)
	cFTOL.SetShowStatusBar(false)
	cFTOL.SetFilteringEnabled(false)
	cFTOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cFTOL.SetShowHelp(true)
	cFTOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	cFTOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cFTOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cFTOL, constant.MaxChooseFetchTagOptionPopUpWidth)

	m.PopUpModel = &ChooseFetchTagOptionPopUpModel{
		RemoteName:      remoteName,
		FetchOptionList: cFTOL,
	}
}

// ------------------------------------
//
//	Initialize the fetch tag output popup. Creates a soft-wrap viewport sized to
//	a fixed height and 80% of terminal width (minus padding), and a dot spinner.
//	All atomic state flags (IsProcessing, HasError, ProcessSuccess, IsCancelled)
//	are reset to false.
//
// ------------------------------------
func InitFetchTagOutputPopUpModel(m *types.GittiModel) {
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpFetchTagOutputViewportHeight)
	vp.SetWidth(min(constant.MaxFetchTagOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &FetchTagOutputPopUpModel{
		FetchTagOutputViewport: vp,
		Spinner:                s,
		CancelFunc:             nil,
	}

	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)

	m.PopUpModel = popUpModel
}
