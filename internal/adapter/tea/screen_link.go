package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateLink(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right":
			switch m.textInput {
			case "submit":
				m.textInput = "public"
				return m, nil
			case "public":
				m.textInput = "submit"
				return m, nil
			}
		case "tab":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.textinputLinkAddress.Blur()
				m.textinputLinkDescription.Focus()
			case "description":
				m.textInput = "submit"
				m.textinputLinkDescription.Blur()
			case "submit":
				m.textInput = "public"
			case "public":
				m.textInput = "address"
				m.textinputLinkAddress.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "address":
				m.textInput = "submit"
				m.textinputLinkAddress.Blur()
				return m, nil
			case "description":
				if len(m.textinputLinkDescription.Value()) == 0 {
					m.textInput = "address"
					m.textinputLinkAddress.Focus()
					m.textinputLinkDescription.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "description"
				m.textinputLinkDescription.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.textinputLinkAddress.Blur()
				m.textinputLinkDescription.Focus()
				return m, nil
			case "description":
				if len(m.textinputLinkDescription.Value()) == 0 {
					m.textInput = "submit"
					m.textinputLinkDescription.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "address"
				m.textinputLinkAddress.Focus()
				return m, nil
			}
		case "esc":
			return m, func() tea.Msg { return "exit" }
		case "q":
			if m.textInput != "address" && m.textInput != "description" {
				return m, func() tea.Msg { return "exit" }
			}
		case "enter":
			switch m.textInput {
			case "address":
				m.textInput = "description"
				m.textinputLinkAddress.Blur()
				m.textinputLinkDescription.Focus()
				return m, nil
			case "description":
				if len(m.textinputLinkDescription.Value()) == 0 {
					m.textInput = "submit"
					m.textinputLinkDescription.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					m.selectedLink.Address = m.textinputLinkAddress.Value()
					m.selectedLink.Description = m.textinputLinkDescription.Value()
					m.selectedItem.UpdateLink(m.selectedLink)
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return "exit"
				}
			case "public":
				return m, func() tea.Msg {
					m.selectedLink.Public = !m.selectedLink.Public
					m.selectedItem.UpdateLink(m.selectedLink)
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
			m.currentScreen = screenLinks
			return m, m.getItem
		}
	}

	switch m.textInput {
	case "address":
		m.textinputLinkAddress, cmd = m.textinputLinkAddress.Update(msg)
		return m, cmd
	case "description":
		m.textinputLinkDescription, cmd = m.textinputLinkDescription.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewLink() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedLink.ID))
	s.WriteString("Address:\n  ")
	s.WriteString(m.textinputLinkAddress.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputLinkDescription.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ SUBMIT ]] ")
	} else {
		s.WriteString("[ submit ] ")
	}

	if m.textInput == "public" {
		if m.selectedLink.Public {
			s.WriteString("[[ SUBMIT AS PRIVATE ]]\n")
		} else {
			s.WriteString("[[ SUBMIT AS PUBLIC ]]\n")
		}
	} else {
		if m.selectedLink.Public {
			s.WriteString("[ submit as private ]\n")
		} else {
			s.WriteString("[ submit as public ]\n")
		}
	}

	return s.String()
}

func (m *TeaModel) resetItemLink() {
	m.textInput = "address"

	m.textinputLinkAddress = textinput.New()
	m.textinputLinkAddress.Placeholder = "link"
	m.textinputLinkAddress.Prompt = ""
	m.textinputLinkAddress.Focus()
	m.textinputLinkAddress.Width = 80
	m.textinputLinkAddress.CharLimit = 1000
	m.textinputLinkAddress.SetValue(m.selectedLink.Address)

	m.textinputLinkDescription = textarea.New()
	m.textinputLinkDescription.Placeholder = "link description"
	m.textinputLinkDescription.Blur()
	m.textinputLinkDescription.Prompt = "  "
	m.textinputLinkDescription.ShowLineNumbers = false
	m.textinputLinkDescription.SetHeight(4)
	m.textinputLinkDescription.SetWidth(80)
	m.textinputLinkDescription.CharLimit = 1000
	m.textinputLinkDescription.SetValue(m.selectedLink.Description)
}
