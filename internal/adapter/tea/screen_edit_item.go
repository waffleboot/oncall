package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateEditItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editItemMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenAllItems
			return m, nil
		}
	}

	return m, nil
}

func (m *TeaModel) viewEditItem() string {
	return fmt.Sprintf("#%d\n", m.items[m.selectedItem].ID)
}

func (m *TeaModel) resetEditItemMenu() {
	item := m.items[m.selectedItem]

	m.editItemMenu.ResetMenu()

	m.editItemMenu.AddGroup("exit")

	if !item.IsClosed() {
		m.editItemMenu.AddGroup("edit_type")
	}

	m.editItemMenu.AddGroup("nodes")
	m.editItemMenu.AddGroup("vms")
	m.editItemMenu.AddGroup("notes")
	m.editItemMenu.AddGroup("links")
	m.editItemMenu.AddDelimiter()

	if item.IsActive() {
		m.editItemMenu.AddGroup("sleep")
	}

	if item.IsSleep() {
		m.editItemMenu.AddGroup("awake")
	}

	if !item.IsClosed() {
		m.editItemMenu.AddGroup("close")
	}

	m.editItemMenu.AddDelimiter()
	m.editItemMenu.AddGroup("delete")
}
