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
		case "tab":
			switch m.textInput {
			case "name":
				m.textInput = "submit"
				m.textinputNodeName.Blur()
			case "submit":
				m.textInput = "name"
				m.textinputNodeName.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "name":
				m.textInput = "submit"
				m.textinputNodeName.Blur()
				return m, nil
			case "submit":
				m.textInput = "name"
				m.textinputNodeName.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "name":
				m.textInput = "submit"
				m.textinputNodeName.Blur()
				return m, nil
			case "submit":
				m.textInput = "name"
				m.textinputNodeName.Focus()
				return m, nil
			}
		case "esc":
			m.currentScreen = screenNodes
			return m, m.getItem
		case "q":
			if m.textInput != "name" {
				m.currentScreen = screenNodes
				return m, m.getItem
			}
		case "enter":
			switch m.textInput {
			case "name":
				m.textInput = "submit"
				m.textinputNodeName.Blur()
				return m, nil
			case "submit":
				return m, func() tea.Msg {
					m.selectedNode.Name = m.textinputNodeName.Value()
					m.selectedItem.UpdateNode(m.selectedNode)
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		}
	case model.Item:
		m.currentScreen = screenNodes
		return m, m.getItem
	}

	switch m.textInput {
	case "name":
		m.textinputNodeName, cmd = m.textinputNodeName.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewNode() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedNode.ID))
	s.WriteString("Name:\n  ")
	s.WriteString(m.textinputNodeName.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ SUBMIT ]]\n")
	} else {
		s.WriteString("[ submit ]\n")
	}

	return s.String()
}

func (m *TeaModel) resetNode() {
	m.textInput = "name"

	m.textinputNodeName = textinput.New()
	m.textinputNodeName.Placeholder = "name"
	m.textinputNodeName.Prompt = ""
	m.textinputNodeName.Focus()
	m.textinputNodeName.Width = 80
	m.textinputNodeName.CharLimit = 1000
	m.textinputNodeName.SetValue(m.selectedNode.Name)
}
