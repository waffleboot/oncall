package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateItemTitle(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.textInput {
			case "title":
				m.textInput = "description"
				m.textinputItemTitle.Blur()
				m.textinputItemDescription.Focus()
			case "description":
				m.textInput = "submit"
				m.textinputItemDescription.Blur()
			case "submit":
				m.textInput = "title"
				m.textinputItemTitle.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "title":
				m.textInput = "submit"
				m.textinputItemTitle.Blur()
				return m, nil
			case "submit":
				m.textInput = "description"
				m.textinputItemDescription.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "title":
				m.textInput = "description"
				m.textinputItemTitle.Blur()
				m.textinputItemDescription.Focus()
				return m, nil
			case "submit":
				m.textInput = "title"
				m.textinputItemTitle.Focus()
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
			case "title":
				m.textInput = "description"
				m.textinputItemTitle.Blur()
				m.textinputItemDescription.Focus()
				return m, nil
			case "description":
				if len(m.textinputItemDescription.Value()) == 0 {
					m.textInput = "submit"
					m.textinputItemDescription.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					m.selectedItem.Title = m.textinputItemTitle.Value()
					m.selectedItem.Description = m.textinputItemDescription.Value()
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
			m.currentScreen = screenItem
			return m, m.getItem
		}
	}

	switch m.textInput {
	case "title":
		m.textinputItemTitle, cmd = m.textinputItemTitle.Update(msg)
		return m, cmd
	case "description":
		m.textinputItemDescription, cmd = m.textinputItemDescription.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewItemTitle() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("#%d - %s\n", m.selectedItem.ID, m.selectedItem.Type.String()))
	s.WriteString("Title:\n  ")
	s.WriteString(m.textinputItemTitle.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputItemDescription.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ Submit ]]\n")
	} else {
		s.WriteString("[ Submit ]\n")
	}

	return s.String()
}

func (m *TeaModel) resetItemTitle() {
	m.textInput = "title"

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
}
