package interaction

import (
	tea "charm.land/bubbletea/v2"
	"github.com/gohyuhan/gitti/tui/constant"
	"github.com/gohyuhan/gitti/tui/interaction/handler"
	"github.com/gohyuhan/gitti/tui/types"
)

func GittiMouseInteraction(msg tea.MouseMsg, m *types.GittiModel) (*types.GittiModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "wheelleft":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
				m.DetailPanelTwoViewport.ScrollLeft(1)
			} else {
				m.DetailPanelViewport.ScrollLeft(1)
			}
		}

	case "wheelright":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
				m.DetailPanelTwoViewport.ScrollRight(1)
			} else {
				m.DetailPanelViewport.ScrollRight(1)
			}
		}

	case "wheelup":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
				m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
			} else {
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			}
			return m, cmd
		} else {
			return handler.UpDownMouseMsgUpdateForPopUp(msg, m)
		}

	case "wheeldown":
		if !m.ShowPopUp.Load() {
			if m.CurrentSelectedComponent == constant.DetailComponentTwo {
				m.DetailPanelTwoViewport, cmd = m.DetailPanelTwoViewport.Update(msg)
			} else {
				m.DetailPanelViewport, cmd = m.DetailPanelViewport.Update(msg)
			}
			return m, cmd
		} else {
			return handler.UpDownMouseMsgUpdateForPopUp(msg, m)
		}
	}
	return m, nil
}
