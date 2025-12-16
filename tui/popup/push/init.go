package push

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	"github.com/gohyuhan/gitti/api/git"
	"github.com/gohyuhan/gitti/i18n"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/style"
	"github.com/gohyuhan/gitti/tui/types"
	"github.com/gohyuhan/gitti/tui/utils"
)

// init the popup model for push output log
func InitGitRemotePushPopUpModel(m *types.GittiModel) {
	// for git push output viewport,
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpGitRemotePushOutputViewportHeight)
	vp.SetWidth(min(constant.MaxGitRemotePushPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitRemotePushPopUpModel{
		GitRemotePushOutputViewport: vp,
		Spinner:                     s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}

// init the popup model for choosing push type
func InitChoosePushTypePopUpModel(m *types.GittiModel, remoteName string) {
	pushTypeOption := []GitPushOptionItem{
		{
			Name:     i18n.LANGUAGEMAPPING.NormalPush,
			Info:     "git push",
			PushType: git.PUSH,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushSafe,
			Info:     "git push --force-with-lease",
			PushType: git.FORCEPUSHSAFE,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.ForcePushDangerous,
			Info:     "git push --force",
			PushType: git.FORCEPUSHDANGEROUS,
		},
	}

	items := make([]list.Item, 0, len(pushTypeOption))
	for _, pushOption := range pushTypeOption {
		items = append(items, GitPushOptionItem(pushOption))
	}
	width := (min(constant.MaxChoosePushTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	pOL := list.New(items, GitPushOptionDelegate{}, width, constant.PopUpChoosePushTypeHeight)
	pOL.SetShowPagination(false)
	pOL.SetShowStatusBar(false)
	pOL.SetFilteringEnabled(false)
	pOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	pOL.SetShowHelp(true)
	pOL.KeyMap = list.KeyMap{} // Clear default keybindings to hide them
	pOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	pOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &pOL, constant.MaxChoosePushTypePopUpWidth)

	m.PopUpModel = &ChoosePushTypePopUpModel{
		PushOptionList: pOL,
		RemoteName:     remoteName,
	}
}
