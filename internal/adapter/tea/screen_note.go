package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
	"github.com/waffleboot/oncall/pkg/tea/tabs"
)

func (m *TeaModel) updateNote(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
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
			if m.submitNote.Focused() {
				return m.exitScreen()
			}
		}
	case button.PressedMsg:
		return m.runAndExitScreen(func() error {
			m.selectedNote.Text = m.textinputNote.Value()
			m.selectedItem.UpdateNote(m.selectedNote)
			if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
				return fmt.Errorf("update item: %w", err)
			}
			return nil
		})
	case string:
		if msg == "exit" {
			m.currentScreen = screenNotes
			return m, m.getItem
		}
	}

	switch {
	case m.textinputNote.Focused():
		m.textinputNote, cmd = m.textinputNote.Update(msg)
		return m, cmd
	case m.submitNote.Focused():
		m.submitNote, cmd = m.submitNote.Update(msg)
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
	s.WriteString(m.submitNote.View())
	s.WriteString("\n")

	return s.String()
}

func (m *TeaModel) resetNote() {
	m.textinputNote = textarea.New()
	m.textinputNote.Placeholder = "note"
	m.textinputNote.Focus()
	m.textinputNote.ShowLineNumbers = false
	m.textinputNote.SetHeight(16)
	m.textinputNote.SetWidth(80)
	m.textinputNote.CharLimit = 1000
	m.textinputNote.SetValue(m.selectedNote.Text)

	m.submitNote = button.New("submit")
	m.submitNote.Blur()

	m.tabs = tabs.New()
	m.tabs.Items = []tabs.Item{
		&m.textinputNote,
		&m.submitNote,
	}
	m.tabs.CanUp = func(tab int) bool {
		if tab == 0 && m.textinputNote.Value() != "" {
			return false
		}
		return true
	}
	m.tabs.CanDown = m.tabs.CanUp
}
