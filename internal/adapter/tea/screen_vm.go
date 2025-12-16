package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateVM(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.textInput {
			case "name":
				m.textInput = "node"
				m.textinputVmName.Blur()
				m.textinputVmNode.Focus()
			case "node":
				m.textInput = "description"
				m.textinputVmNode.Blur()
				m.textinputVmDescription.Focus()
			case "description":
				m.textInput = "submit"
				m.textinputVmDescription.Blur()
			case "submit":
				m.textInput = "name"
				m.textinputVmName.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "name":
				m.textInput = "submit"
				m.textinputVmName.Blur()
				return m, nil
			case "node":
				m.textInput = "name"
				m.textinputVmNode.Blur()
				m.textinputVmName.Focus()
				return m, nil
			case "description":
				if len(m.textinputVmDescription.Value()) == 0 {
					m.textInput = "node"
					m.textinputVmNode.Focus()
					m.textinputVmDescription.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "description"
				m.textinputVmDescription.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "name":
				m.textInput = "node"
				m.textinputVmName.Blur()
				m.textinputVmNode.Focus()
				return m, nil
			case "node":
				m.textInput = "description"
				m.textinputVmNode.Blur()
				m.textinputVmDescription.Focus()
				return m, nil
			case "description":
				if len(m.textinputVmDescription.Value()) == 0 {
					m.textInput = "submit"
					m.textinputVmDescription.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "name"
				m.textinputVmName.Focus()
				return m, nil
			}
		case "esc":
			m.currentScreen = screenVMs
			return m, m.getItem
		case "q":
			if m.textInput != "name" && m.textInput != "node" && m.textInput != "description" {
				m.currentScreen = screenVMs
				return m, m.getItem
			}
		case "enter":
			switch m.textInput {
			case "name":
				m.textInput = "node"
				m.textinputVmName.Blur()
				m.textinputVmNode.Focus()
				return m, nil
			case "node":
				m.textInput = "description"
				m.textinputVmNode.Blur()
				m.textinputVmDescription.Focus()
				return m, nil
			case "description":
				if len(m.textinputVmDescription.Value()) == 0 {
					m.textInput = "submit"
					m.textinputVmDescription.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					m.selectedVM.Name = m.textinputVmName.Value()
					m.selectedVM.Node = m.textinputVmNode.Value()
					m.selectedVM.Description = m.textinputVmDescription.Value()
					m.selectedItem.UpdateVM(m.selectedVM)
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		}
	case model.Item:
		m.currentScreen = screenVMs
		return m, m.getItem
	}

	switch m.textInput {
	case "name":
		m.textinputVmName, cmd = m.textinputVmName.Update(msg)
		return m, cmd
	case "node":
		m.textinputVmNode, cmd = m.textinputVmNode.Update(msg)
		return m, cmd
	case "description":
		m.textinputVmDescription, cmd = m.textinputVmDescription.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewVM() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedVM.ID))
	s.WriteString("Name:\n  ")
	s.WriteString(m.textinputVmName.View())
	s.WriteString("\nNode:\n  ")
	s.WriteString(m.textinputVmNode.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputVmDescription.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ SUBMIT ]]\n")
	} else {
		s.WriteString("[ submit ]\n")
	}

	return s.String()
}

func (m *TeaModel) resetVM() {
	m.textInput = "name"

	m.textinputVmName = textinput.New()
	m.textinputVmName.Placeholder = "name"
	m.textinputVmName.Prompt = ""
	m.textinputVmName.Focus()
	m.textinputVmName.Width = 80
	m.textinputVmName.CharLimit = 1000
	m.textinputVmName.SetValue(m.selectedVM.Name)

	m.textinputVmNode = textinput.New()
	m.textinputVmNode.Placeholder = "node"
	m.textinputVmNode.Prompt = ""
	m.textinputVmNode.Blur()
	m.textinputVmNode.Width = 80
	m.textinputVmNode.CharLimit = 1000
	m.textinputVmNode.SetValue(m.selectedVM.Node)

	m.textinputVmDescription = textarea.New()
	m.textinputVmDescription.Placeholder = "link description"
	m.textinputVmDescription.Blur()
	m.textinputVmDescription.Prompt = "  "
	m.textinputVmDescription.ShowLineNumbers = false
	m.textinputVmDescription.SetHeight(4)
	m.textinputVmDescription.SetWidth(80)
	m.textinputVmDescription.CharLimit = 1000
	m.textinputVmDescription.SetValue(m.selectedVM.Description)
}
