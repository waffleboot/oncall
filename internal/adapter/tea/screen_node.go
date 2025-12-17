package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateNode(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "enter":
			return m.runAndExitScreen(func() error {
				m.selectedNode.Name = m.textinputNode.Value()
				m.selectedItem.UpdateNode(m.selectedNode)
				if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				}
				return nil
			})
		}
	case string:
		if msg == "exit" {
			m.currentScreen = screenNodes
			return m, m.getItem
		}
	}

	m.textinputNode, cmd = m.textinputNode.Update(msg)

	return m, cmd
}

func (m *TeaModel) viewNode() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedNode.ID))
	s.WriteString("Node:\n  ")
	s.WriteString(m.textinputNode.View())
	s.WriteString("\n")
	return s.String()
}

func (m *TeaModel) resetNode() {
	m.textinputNode = textinput.New()
	m.textinputNode.Placeholder = "node"
	m.textinputNode.Prompt = ""
	m.textinputNode.Focus()
	m.textinputNode.Width = 80
	m.textinputNode.CharLimit = 1000
	m.textinputNode.SetValue(m.selectedNode.Name)
}
