package tea

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateNode(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.currentScreen = screenNodes
			return m, m.getItem
		case "enter":
			return m, func() tea.Msg {
				if m.selectedNode.ID == 0 {
					s := bufio.NewScanner(strings.NewReader(m.textinputNodeNames.Value()))
					for s.Scan() {
						t := strings.TrimSpace(s.Text())
						if t != "" {
							node := m.selectedItem.CreateNode()
							node.Name = t
							m.selectedItem.UpdateNode(node)
						}
					}
					if err := s.Err(); err != nil {
						return fmt.Errorf("scan: %w", err)
					}
				} else {
					m.selectedNode.Name = m.textinputNodeName.Value()
					m.selectedItem.UpdateNode(m.selectedNode)
				}
				if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				}
				return m.getItem()
			}
		}
	case model.Item:
		m.currentScreen = screenNodes
		return m, m.getItem
	}

	if m.selectedNode.ID == 0 {
		m.textinputNodeNames, cmd = m.textinputNodeNames.Update(msg)
	} else {
		m.textinputNodeName, cmd = m.textinputNodeName.Update(msg)
	}
	return m, cmd
}

func (m *TeaModel) viewNode() string {
	var s strings.Builder
	if m.selectedNode.ID == 0 {
		s.WriteString("Nodes:\n  ")
		s.WriteString(m.textinputNodeNames.View())
		s.WriteString("\n")
	} else {
		s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedNode.ID))
		s.WriteString("Node:\n  ")
		s.WriteString(m.textinputNodeName.View())
		s.WriteString("\n")
	}
	return s.String()
}

func (m *TeaModel) resetNode() {
	if m.selectedNode.ID == 0 {
		m.textinputNodeNames = textarea.New()
		m.textinputNodeNames.Placeholder = "nodes"
		m.textinputNodeNames.Focus()
		m.textinputNodeNames.Prompt = "  "
		m.textinputNodeNames.ShowLineNumbers = false
		m.textinputNodeNames.SetHeight(4)
		m.textinputNodeNames.SetWidth(80)
		m.textinputNodeNames.CharLimit = 1000
		m.textinputNodeNames.SetValue(m.selectedNode.Name)
	} else {
		m.textinputNodeName = textinput.New()
		m.textinputNodeName.Placeholder = "node"
		m.textinputNodeName.Prompt = ""
		m.textinputNodeName.Focus()
		m.textinputNodeName.Width = 80
		m.textinputNodeName.CharLimit = 1000
		m.textinputNodeName.SetValue(m.selectedNode.Name)
	}
}
