package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateNodes(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuNodes.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
			return m, m.getItem
		case "d":
			if g, p := m.menuNodes.GetGroup(); g == "nodes" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteNode(m.nodes[p])
					if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			m.currentScreen = screenNewNodes
			m.resetNewNodes()
			return m, nil
		case "enter", " ":
			switch g, p := m.menuNodes.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				m.currentScreen = screenNewNodes
				m.resetNewNodes()
				return m, nil
			case "nodes":
				return m, func() tea.Msg { return m.nodes[p] }
			}
		}
	case model.Item:
		m.selectedItem = msg
		m.resetNodes("")
	case model.Node:
		m.selectedNode = msg
		m.currentScreen = screenNode
		m.resetNode()
	}

	return m, nil
}

func (m *TeaModel) viewNodes() string {
	return m.menuNodes.View()
}

func (m *TeaModel) resetNodes(toGroup string) {
	m.nodes = m.selectedItem.ActiveNodes()
	m.menuNodes.ResetMenu()
	m.menuNodes.AddGroup("exit")
	m.menuNodes.AddGroup("new")
	m.menuNodes.AddGroupWithItems("nodes", len(m.nodes))
	if toGroup != "" {
		m.menuNodes.JumpToGroup(toGroup)
	} else {
		m.menuNodes.AdjustCursor()
	}
}
