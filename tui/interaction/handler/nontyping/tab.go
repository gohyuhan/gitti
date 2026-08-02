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
//	Handle Tab key interaction.
//	Responsibility: Forward navigation.
//	- No popup: Cycles the primary active focus to the *next* main component panel in the UI layout
//	  (e.g., from Branches to Modified Files, then to Commit Log, etc.).
//	- Merge Popup: Shifts section focus from the branch option list to the selected branch list.
//
// ------------------------------------
func handleNonTypingTabKeyBindingInteraction(m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	if !m.ShowPopUp.Load() {
		nextNavigation := m.CurrentSelectedComponentIndex + 1
		if nextNavigation < len(constant.ComponentPanelNavigationList) {
			m.CurrentSelectedComponentIndex = nextNavigation
			m.CurrentSelectedComponent = constant.ComponentPanelNavigationList[nextNavigation]
			m.DetailPanelParentComponent = ""
			layout.LeftPanelDynamicResize(m)
			services.FetchDetailComponentPanelInfoService(m, true)
		}
	} else {
		switch m.PopUpType {
		case constant.ChooseBranchOptionForMergePopUp:
			popUp, ok := m.PopUpModel.(*branchPopUp.ChooseBranchOptionForMergePopUpModel)
			if ok {
				if popUp.BranchOptionSectionSelected.Load() {
					popUp.BranchOptionSectionSelected.Store(false)
					popUp.SelectedBranchSectionSelected.Store(true)
				}
			}
		}
	}
	return m, nil
}
