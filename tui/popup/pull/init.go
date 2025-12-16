package pull

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

func InitChooseGitPullTypePopUp(m *types.GittiModel) {
	pullTypeOption := []GitPullTypeOptionItem{
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullOption,
			Info:     "git pull --no-edit",
			PullType: git.GITPULL,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullRebaseOption,
			Info:     "git pull --rebase --autostash --no-edit",
			PullType: git.GITPULLREBASE,
		},
		{
			Name:     i18n.LANGUAGEMAPPING.GitPullMergeOption,
			Info:     "git pull --no-rebase --no-edit",
			PullType: git.GITPULLMERGE,
		},
	}

	items := make([]list.Item, 0, len(pullTypeOption))
	for _, pullOption := range pullTypeOption {
		items = append(items, GitPullTypeOptionItem(pullOption))
	}

	width := (min(constant.MaxChooseGitPullTypePopUpWidth, int(float64(m.Width)*0.8)) - 4)
	cGPOL := list.New(items, GitPullTypeOptionDelegate{}, width, constant.PopUpChooseGitPullTypeHeight)
	cGPOL.SetShowPagination(false)
	cGPOL.SetShowStatusBar(false)
	cGPOL.SetFilteringEnabled(false)
	cGPOL.SetShowTitle(false)

	// Custom Help Model for Count Display
	cGPOL.SetShowHelp(true)
	cGPOL.KeyMap = list.KeyMap{}
	cGPOL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	cGPOL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &cGPOL, constant.MaxChooseGitPullTypePopUpWidth)

	popUpModel := &ChooseGitPullTypePopUpModel{
		PullTypeOptionList: cGPOL,
	}

	m.PopUpModel = popUpModel
}

func InitGitPullOutputPopUpModel(m *types.GittiModel) {
	// for git pull output viewport
	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(constant.PopUpSwitchBranchOutputViewPortHeight)
	vp.SetWidth(min(constant.MaxSwitchBranchOutputPopUpWidth, int(float64(m.Width)*0.8)) - 4)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style.SpinnerStyle

	popUpModel := &GitPullOutputPopUpModel{
		GitPullOutputViewport: vp,
		Spinner:               s,
	}
	popUpModel.IsProcessing.Store(false)
	popUpModel.HasError.Store(false)
	popUpModel.ProcessSuccess.Store(false)
	popUpModel.IsCancelled.Store(false)
	m.PopUpModel = popUpModel
}
