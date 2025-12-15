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
				m.itemTitle.Blur()
				m.itemDescription.Focus()
			case "description":
				m.textInput = "submit"
				m.itemDescription.Blur()
			case "submit":
				m.textInput = "title"
				m.itemTitle.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "title":
				m.textInput = "submit"
				m.itemTitle.Blur()
				return m, nil
			case "submit":
				m.textInput = "description"
				m.itemDescription.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "title":
				m.textInput = "description"
				m.itemTitle.Blur()
				m.itemDescription.Focus()
				return m, nil
			case "submit":
				m.textInput = "title"
				m.itemTitle.Focus()
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
				m.itemTitle.Blur()
				m.itemDescription.Focus()
				return m, nil
			case "description":
				if len(m.itemDescription.Value()) == 0 {
					m.textInput = "submit"
					m.itemDescription.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					m.selectedItem.Title = m.itemTitle.Value()
					m.selectedItem.Description = m.itemDescription.Value()
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
			m.currentScreen = screenEditItem
			return m, m.getItem
		}
	}

	switch m.textInput {
	case "title":
		m.itemTitle, cmd = m.itemTitle.Update(msg)
		return m, cmd
	case "description":
		m.itemDescription, cmd = m.itemDescription.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewItemTitle() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("#%d - %s\n", m.selectedItem.ID, m.selectedItem.Type.String()))
	s.WriteString("Title:\n  ")
	s.WriteString(m.itemTitle.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.itemDescription.View())
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

	m.itemTitle = textinput.New()
	m.itemTitle.Placeholder = "title"
	m.itemTitle.Prompt = ""
	m.itemTitle.Focus()
	m.itemTitle.Width = 80
	m.itemTitle.CharLimit = 1000
	m.itemTitle.SetValue(m.selectedItem.Title)

	m.itemDescription = textarea.New()
	m.itemDescription.Placeholder = "description"
	m.itemDescription.Blur()
	m.itemDescription.Prompt = "  "
	m.itemDescription.ShowLineNumbers = false
	m.itemDescription.SetHeight(4)
	m.itemDescription.SetWidth(80)
	m.itemDescription.CharLimit = 1000
	m.itemDescription.SetValue(m.selectedItem.Description)
}
