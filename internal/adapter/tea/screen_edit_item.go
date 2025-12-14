package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
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
	var state string

	item := m.items[m.selectedItem]

	switch {
	case item.IsSleep():
		state = " в ожидании"
	case item.IsClosed():
		switch item.Type {
		case model.ItemTypeAsk:
			state = " закрыто"
		default:
			state = " закрыт"
		}
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("  #%d %s%s\n\n", item.ID, item.Type, state))
	s.WriteString(m.editItemMenu.GenerateMenu())

	return s.String()
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
