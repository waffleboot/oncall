package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateNotes(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuNotes.Update(msg) {
		return m, nil
	}

	newNote := func() tea.Msg {
		return m.selectedItem.CreateNote()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
			return m, m.getItem
		case "d":
			if g, p := m.menuNotes.GetGroup(); g == "notes" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteNote(m.notes[p])
					if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			return m, newNote
		case "p":
			if g, p := m.menuNotes.GetGroup(); g == "notes" {
				return m, func() tea.Msg {
					note := m.notes[p]
					note.Public = !note.Public
					m.selectedItem.UpdateNote(note)
					if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "enter", " ":
			switch g, p := m.menuNotes.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				return m, newNote
			case "notes":
				return m, func() tea.Msg { return m.notes[p] }
			}
		}
	case model.Item:
		m.selectedItem = msg
		m.resetNotes("")
	case model.Note:
		m.selectedNote = msg
		m.currentScreen = screenNote
		m.resetNote()
	}

	return m, nil
}

func (m *TeaModel) viewNotes() string {
	return m.menuNotes.View()
}

func (m *TeaModel) resetNotes(toGroup string) {
	m.notes = m.selectedItem.ActiveNotes()
	m.menuNotes.ResetMenu()
	m.menuNotes.AddGroup("exit")
	m.menuNotes.AddGroup("new")
	m.menuNotes.AddGroupWithItems("notes", len(m.notes))
	if toGroup != "" {
		m.menuNotes.JumpToGroup(toGroup)
	} else {
		m.menuNotes.AdjustCursor()
	}
}
