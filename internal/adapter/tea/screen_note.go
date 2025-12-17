package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateNote(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.textInput {
			case "text":
				m.textInput = "submit"
				m.textinputNote.Blur()
			case "submit":
				m.textInput = "text"
				m.textinputNote.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "text":
				if len(m.textinputNote.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNote.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "text"
				m.textinputNote.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "text":
				if len(m.textinputNote.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNote.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "text"
				m.textinputNote.Focus()
				return m, nil
			}
		case "esc":
			m.currentScreen = screenNotes
			return m, m.getItem
		case "q":
			if m.textInput != "text" {
				m.currentScreen = screenNotes
				return m, m.getItem
			}
		case "enter":
			switch m.textInput {
			case "text":
				if len(m.textinputNote.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNote.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
					note := m.textinputNote.Value()
					if strings.TrimSpace(note) != "" {
						m.selectedNote.Text = note
						m.selectedItem.UpdateNote(m.selectedNote)
						if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
							return fmt.Errorf("update item: %w", err)
						}
					}
					return m.getItem()
				}
			}
		}
	case model.Item:
		m.currentScreen = screenNotes
		return m, m.getItem
	}

	switch m.textInput {
	case "text":
		m.textinputNote, cmd = m.textinputNote.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewNote() string {
	var s strings.Builder

	if m.selectedNote.Exists() {
		s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedNote.ID))
	}
	s.WriteString("Note:\n")
	s.WriteString(m.textinputNote.View())
	s.WriteString("\n")

	if m.textInput == "submit" {
		s.WriteString("[[ SUBMIT ]]\n")
	} else {
		s.WriteString("[ submit ]\n")
	}

	return s.String()
}

func (m *TeaModel) resetNote() {
	m.textInput = "text"
	m.textinputNote = textarea.New()
	m.textinputNote.Placeholder = "note"
	m.textinputNote.Focus()
	m.textinputNote.ShowLineNumbers = false
	m.textinputNote.SetHeight(4)
	m.textinputNote.SetWidth(80)
	m.textinputNote.CharLimit = 1000
	m.textinputNote.SetValue(m.selectedNote.Text)
}
