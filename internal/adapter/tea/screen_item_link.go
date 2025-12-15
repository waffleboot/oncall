package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateItemLink(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.linkAddress.Blur()
				m.linkDescription.Focus()
			case "description":
				m.textInput = "submit"
				m.linkDescription.Blur()
			case "submit":
				m.textInput = "address"
				m.linkAddress.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "address":
				m.textInput = "submit"
				m.linkAddress.Blur()
				return m, nil
			case "submit":
				m.textInput = "description"
				m.linkDescription.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.linkAddress.Blur()
				m.linkDescription.Focus()
				return m, nil
			case "submit":
				m.textInput = "address"
				m.linkAddress.Focus()
				return m, nil
			}
		case "esc":
			return m, func() tea.Msg { return "exit" }
		case "q":
			if m.textInput == "submit" {
				return m, func() tea.Msg { return "exit" }
			}
		case "enter":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.linkAddress.Blur()
				m.linkDescription.Focus()
				return m, nil
			case "description":
				if len(m.linkDescription.Value()) == 0 {
					m.textInput = "submit"
					m.linkDescription.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					m.selectedLink.Address = m.linkAddress.Value()
					m.selectedLink.Description.SetValue(m.linkDescription.Value())
					m.selectedItem.UpdateItemLink(m.selectedLink)
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return "exit"
				}
			}
		}
	case string:
		switch msg {
		case "exit":
			m.currentScreen = screenItemLinks
			return m, m.getItem
		}
	}

	switch m.textInput {
	case "address":
		m.linkAddress, cmd = m.linkAddress.Update(msg)
		return m, cmd
	case "description":
		m.linkDescription, cmd = m.linkDescription.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewItemLink() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedLink.ID))
	s.WriteString("Address:\n  ")
	s.WriteString(m.linkAddress.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.linkDescription.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ Submit ]]\n")
	} else {
		s.WriteString("[ Submit ]\n")
	}

	return s.String()
}

func (m *TeaModel) resetItemLink() {
	m.textInput = "address"

	m.linkAddress = textinput.New()
	m.linkAddress.Placeholder = "link"
	m.linkAddress.Prompt = ""
	m.linkAddress.Focus()
	m.linkAddress.Width = 80
	m.linkAddress.CharLimit = 1000
	m.linkAddress.SetValue(m.selectedLink.Address)

	m.linkDescription = textarea.New()
	m.linkDescription.Placeholder = "link description"
	m.linkDescription.Blur()
	m.linkDescription.Prompt = "  "
	m.linkDescription.ShowLineNumbers = false
	m.linkDescription.SetHeight(4)
	m.linkDescription.SetWidth(80)
	m.linkDescription.CharLimit = 1000
	m.linkDescription.SetValue(m.selectedLink.Description.Value())
}
