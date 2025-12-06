package remote

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/types"
)

// init the popup model for prompting user to add remote origin
func InitAddRemotePromptPopUpModel(m *types.GittiModel, noInitialRemote bool) {
	RemoteNameTextInput := textinput.New()
	if noInitialRemote {
		RemoteNameTextInput.SetValue("origin")
	}
	RemoteNameTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteNamePlaceHolder
	RemoteNameTextInput.Focus()
	RemoteNameTextInput.SetVirtualCursor(true)

	RemoteUrlTextInput := textinput.New()
	RemoteUrlTextInput.Placeholder = i18n.LANGUAGEMAPPING.AddRemotePopUpRemoteUrlPlaceHolder
	RemoteUrlTextInput.Blur()
	RemoteUrlTextInput.SetVirtualCursor(true)

	// for git add remote output viewport, we will not have any interaction for it as usually it will be a one line for error log or also for our custom success message
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpAddRemoteOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxAddRemotePromptPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	popUpModel := &AddRemotePromptPopUpModel{
		RemoteNameTextInput:     RemoteNameTextInput,
		RemoteUrlTextInput:      RemoteUrlTextInput,
		TotalInputCount:         2,
		CurrentActiveInputIndex: 1,
		AddRemoteOutputViewport: vp,
		NoInitialRemote:         noInitialRemote,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model to choose remote to push to
func InitGitRemotePushChooseRemotePopUpModel(m *types.GittiModel, remoteList []git.GitRemoteInfo) {
	items := make([]list.Item, 0, len(remoteList))
	for _, remote := range remoteList {
		items = append(items, GitRemoteItem(remote))
	}
	width := (min(constant.MaxChooseRemotePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	rL := list.New(items, GitRemoteItemDelegate{}, width, constant.PopUpChooseRemoteHeight)
	rL.SetShowPagination(false)
	rL.SetShowStatusBar(false)
	rL.SetFilteringEnabled(false)
	rL.SetShowHelp(false)
	rL.SetShowTitle(false)

	m.PopUpModel = &ChooseRemotePopUpModel{
		RemoteList: rL,
	}
}
