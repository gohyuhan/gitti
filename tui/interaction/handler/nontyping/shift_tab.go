package nontyping

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/layout"
	branchPopUp "github.com/gohyuhan/gitti/tui/popup/branch"
	"github.com/gohyuhan/gitti/tui/services"
	"github.com/gohyuhan/gitti/tui/types"
)

// ------------------------------------
//
//	Handle Shift+Tab key interaction.
//	Responsibility: Backward navigation.
//	- No popup: Cycles the primary active focus to the *previous* main component panel in the UI layout.
//	- Merge Popup: Shifts section focus from the selected branch list back to the branch option list.
//
// ------------------------------------
func handleNonTypingShiftTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		previousNavigation := m.CurrentSelectedComponentIndex - 1
		if previousNavigation >= 0 {
			m.CurrentSelectedComponentIndex = previousNavigation
			m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[previousNavigation]
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseBranchOptionForMergePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				if popUp.SelectedBranchSectionSelected.Load() {
					popUp.BranchOptionSectionSelected.Store(true)
					popUp.SelectedBranchSectionSelected.Store(false)
				}
			}
		}
	}
	return m, nil
}
