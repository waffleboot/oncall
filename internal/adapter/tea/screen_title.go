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

func (m *TeaModel) updateItemTitle(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
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
			if m.submitTitle.Focused() {
				return m.exitScreen()
			}
		case "enter":
			if m.textinputItemTitle.Focused() {
				var ok bool
				if m.tabs, cmd, ok = m.tabs.Next(); ok {
					return m, cmd
				}
			}
		}
	case button.PressedMsg:
		switch msg.Value {
		case "submit":
			return m.runAndExitScreen(func() error {
				m.selectedItem.Title = m.textinputItemTitle.Value()
				m.selectedItem.Description = m.textinputItemDescription.Value()
				if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				}
				return nil
			})
		}
	case string:
		if msg == "exit" {
			m.currentScreen = screenItem
			return m, m.getItem
		}
	}

	switch {
	case m.submitTitle.Focused():
		m.submitTitle, cmd = m.submitTitle.Update(msg)
		return m, cmd
	case m.textinputItemTitle.Focused():
		m.textinputItemTitle, cmd = m.textinputItemTitle.Update(msg)
		return m, cmd
	case m.textinputItemDescription.Focused():
		m.textinputItemDescription, cmd = m.textinputItemDescription.Update(msg)
		return m, cmd
	case m.submitTitle.Focused():
		m.submitTitle, cmd = m.submitTitle.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewTitle() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("#%d - %s\n", m.selectedItem.Num, m.selectedItem.Type.String()))
	s.WriteString("Title:\n  ")
	s.WriteString(m.textinputItemTitle.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputItemDescription.View())
	s.WriteString("\n")
	s.WriteString(m.submitTitle.View())
	s.WriteString("\n")

	return s.String()
}

func (m *TeaModel) resetTitle() {
	m.textinputItemTitle = textinput.New()
	m.textinputItemTitle.Placeholder = "title"
	m.textinputItemTitle.Prompt = ""
	m.textinputItemTitle.Focus()
	m.textinputItemTitle.Width = 80
	m.textinputItemTitle.CharLimit = 1000
	m.textinputItemTitle.SetValue(m.selectedItem.Title)

	m.textinputItemDescription = textarea.New()
	m.textinputItemDescription.Placeholder = "description"
	m.textinputItemDescription.Blur()
	m.textinputItemDescription.Prompt = "  "
	m.textinputItemDescription.ShowLineNumbers = false
	m.textinputItemDescription.SetHeight(4)
	m.textinputItemDescription.SetWidth(80)
	m.textinputItemDescription.CharLimit = 1000
	m.textinputItemDescription.SetValue(m.selectedItem.Description)

	m.submitTitle = button.New("submit")
	m.submitTitle.Blur()

	m.tabs = tabs.New()
	m.tabs.Items = []tabs.Item{
		&m.textinputItemTitle,
		&m.textinputItemDescription,
		&m.submitTitle,
	}
	m.tabs.CanUp = func(tab int) bool {
		if tab == 1 && m.textinputItemDescription.Value() != "" {
			return false
		}
		return true
	}
	m.tabs.CanDown = m.tabs.CanUp
}
