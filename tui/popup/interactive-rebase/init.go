package interactiverebase

import (
	"charm.land/bubbles/v2/list"
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
//	Init the popup model for choosing an interactive rebase operation type (fixup/squash, reword, drop)
//
// ------------------------------------
func InitInteractiveRebaseOptionPopUpModel(m *types.GittiModel) {
	interactiveRebaseOption := []InteractiveRebaseOptionItem{
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquash,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFixupSquashDescription,
			InteractiveRebaseType: git.FIXUPSQUASH,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseReword,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.REWORD,
		},
		{
			Name:                  i18n.LANGUAGEMAPPING.InteractiveRebaseDrop,
			Info:                  i18n.LANGUAGEMAPPING.InteractiveRebaseFeatureComingSoon,
			InteractiveRebaseType: git.DROP,
		},
	}

	items := make([]list.Item, 0, len(interactiveRebaseOption))
	for _, interactiveRebaseOption := range interactiveRebaseOption {
		items = append(items, InteractiveRebaseOptionItem(interactiveRebaseOption))
	}

	width := (min(constant.MaxInteractiveRebaseOptionPopUpWidth, int(float64(m.Width)*0.8)) - 4)
	iROL := list.New(items, InteractiveRebaseOptionDelegate{}, width, constant.PopUpInteractiveRebaseOptionHeight)
	iROL.SetShowPagination(false)
	iROL.SetShowStatusBar(false)
	iROL.SetFilteringEnabled(false)
	iROL.SetShowTitle(false)

	// Custom Help Model for Count Display
	iROL.SetShowHelp(true)
	iROL.KeyMap = list.KeyMap{}
	iROL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	iROL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &iROL, constant.MaxInteractiveRebaseOptionPopUpWidth)

	popUpModel := &InteractiveRebaseOptionPopUpModel{
		InteractiveRebaseOptionList: iROL,
	}

	m.PopUpModel = popUpModel
}

// ------------------------------------
//
//	Init the split-pane popup for selecting commits to fixup or squash;
//	commit list is populated asynchronously via a goroutine after init
//
// ------------------------------------
func InitInteractiveRebaseFixupSquashSelectionPopUpModel(m *types.GittiModel) {
	items := make([]list.Item, 0)
	popUpWidth := int(float64(m.Width) * 0.9)
	innerWidth := popUpWidth - 2
	listWidth := int(float64(innerWidth) * 0.65)
	vpWidth := innerWidth - listWidth

	selectedCommitHashMap := make(map[string]git.CommitInfo)
	height := int(float64(m.Height)*0.8) - 2
	iRFSSL := list.New(items, InteractiveRebaseFixupSquashSelectionDelegate{&selectedCommitHashMap}, listWidth-2, height)
	iRFSSL.SetShowPagination(false)
	iRFSSL.SetShowStatusBar(false)
	iRFSSL.SetFilteringEnabled(false)
	iRFSSL.SetShowTitle(false)

	// Custom Help Model for Count Display
	iRFSSL.SetShowHelp(true)
	iRFSSL.KeyMap = list.KeyMap{}
	iRFSSL.Styles.HelpStyle = style.NewStyle.MarginTop(0).MarginBottom(0).PaddingTop(0).PaddingBottom(0)
	iRFSSL.AdditionalShortHelpKeys = utils.PopUpListCounterHelper(m, &iRFSSL, m.Width)

	vp := viewport.New()
	vp.SoftWrap = true
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1
	vp.SetHeight(height)
	vp.SetWidth(vpWidth - 2)

	vp.SetContent(i18n.LANGUAGEMAPPING.InteractiveRebaseFixupMustHaveAtLeastTwoSelectedError)

	popUpModel := &InteractiveRebaseFixupSquashSelectionPopUpModel{
		CommitList:                  iRFSSL,
		CommitFixupSquashViewport:   vp,
		OriginalRetrievedCommitList: []git.CommitInfo{},
		SelectedCommitHashMap:       selectedCommitHashMap,
	}

	go func() {
		commitInfos := m.GitOperations.GitInteractiveRebase.GetCommitInfos()
		newItems := make([]list.Item, 0, len(commitInfos))
		for _, commitInfo := range commitInfos {
			newItems = append(newItems, InteractiveRebaseFixupSquashSelectionItem{
				Hash:        commitInfo.Hash,
				Message:     commitInfo.Message,
				Author:      commitInfo.Author,
				Description: commitInfo.Description,
				Parent:      commitInfo.Parent,
				CommitOrder: commitInfo.CommitOrder,
			})
		}
		popUpModel.OriginalRetrievedCommitList = commitInfos
		popUpModel.CommitList.SetItems(newItems)
	}()

	popUpModel.IsCommitListSelected.Store(true)
	popUpModel.IsCommitFixupSquashViewportSelected.Store(false)

	m.PopUpModel = popUpModel
}
