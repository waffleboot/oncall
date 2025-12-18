package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
	"github.com/waffleboot/oncall/pkg/tea/tabs"
)

func (m *TeaModel) updateLink(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	var ok bool
	if m.tabs, cmd, ok = m.tabs.Update(msg); ok {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "q":
			if m.textinputLinkAddress.Focused() || m.textinputLinkDescription.Focused() {
				break
			}
			return m.exitScreen()
		case "enter":
			if m.textinputLinkAddress.Focused() {
				var ok bool
				if m.tabs, cmd, ok = m.tabs.Next(); ok {
					return m, cmd
				}
			}
		}
	case button.PressedMsg:
		m.selectedLink.Address = m.textinputLinkAddress.Value()
		m.selectedLink.Description = m.textinputLinkDescription.Value()

		switch msg.Value {
		case "submit as public", "submit as private":
			m.selectedLink.Public = msg.Value == "submit as public"
		}

		return m.runAndExitScreen(func() error {
			m.selectedItem.UpdateLink(m.selectedLink)
			if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
				return fmt.Errorf("update item: %w", err)
			}
			return nil
		})
	case string:
		if msg == "exit" {
			m.currentScreen = screenLinks
			return m, m.getItem
		}
	}

	switch {
	case m.textinputLinkAddress.Focused():
		m.textinputLinkAddress, cmd = m.textinputLinkAddress.Update(msg)
		return m, cmd
	case m.textinputLinkDescription.Focused():
		m.textinputLinkDescription, cmd = m.textinputLinkDescription.Update(msg)
		return m, cmd
	case m.submitLink.Focused():
		m.submitLink, cmd = m.submitLink.Update(msg)
		return m, cmd
	case m.submitAsPublicLink.Focused():
		m.submitAsPublicLink, cmd = m.submitAsPublicLink.Update(msg)
		return m, cmd
	case m.submitAsPrivateLink.Focused():
		m.submitAsPrivateLink, cmd = m.submitAsPrivateLink.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewLink() string {
	var s strings.Builder

	if m.selectedLink.Exists() {
		s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedLink.ID))
	}
	s.WriteString("Address:\n  ")
	s.WriteString(m.textinputLinkAddress.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputLinkDescription.View())
	s.WriteString("\n")
	s.WriteString(m.submitLink.View())
	s.WriteString("\n\n")

	if m.tabs.Visible(3) {
		s.WriteString(m.submitAsPublicLink.View())
		s.WriteString("\n")
	}

	if m.tabs.Visible(4) {
		s.WriteString(m.submitAsPrivateLink.View())
		s.WriteString("\n")
	}

	return s.String()
}

func (m *TeaModel) resetItemLink() {
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

	m.submitLink = button.New("submit")
	m.submitAsPublicLink = button.New("submit as public")
	m.submitAsPrivateLink = button.New("submit as private")

	m.tabs = tabs.New()
	m.tabs.Items = []tabs.Item{
		&m.textinputLinkAddress,
		&m.textinputLinkDescription,
		&m.submitLink,
		&m.submitAsPublicLink,
		&m.submitAsPrivateLink,
	}
	m.tabs.CanUp = func(tab int) bool {
		if tab == 1 && m.textinputLinkDescription.Value() != "" {
			return false
		}
		return true
	}
	m.tabs.CanDown = m.tabs.CanUp
	m.tabs.Visible = func(tab int) bool {
		switch {
		case tab == 3 && m.selectedLink.Public:
			return false
		case tab == 4 && !m.selectedLink.Public:
			return false
		default:
			return true
		}
	}
}
