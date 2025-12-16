package tea

import (
	"fmt"
	"strings"

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
				m.selectedNode.Name = m.textinputNodeName.Value()
				m.selectedItem.UpdateNode(m.selectedNode)
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

	m.textinputNodeName, cmd = m.textinputNodeName.Update(msg)
	return m, cmd
}

func (m *TeaModel) viewNode() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedNode.ID))
	s.WriteString("Name:\n  ")
	s.WriteString(m.textinputNodeName.View())
	s.WriteString("\n")
	return s.String()
}

func (m *TeaModel) resetNode() {
	m.textinputNodeName = textinput.New()
	m.textinputNodeName.Placeholder = "name"
	m.textinputNodeName.Prompt = ""
	m.textinputNodeName.Focus()
	m.textinputNodeName.Width = 80
	m.textinputNodeName.CharLimit = 1000
	m.textinputNodeName.SetValue(m.selectedNode.Name)
}
