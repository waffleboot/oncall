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

	item, found := m.getSelectedItem()
	if !found {
		m.currentScreen = screenAllItems
		return m, m.getItems
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenAllItems
		case "enter", " ":
			switch g, _ := m.editItemMenu.GetGroup(); g {
			case "exit":
				m.currentScreen = screenAllItems
			case "sleep":
				return m, func() tea.Msg {
					if item, err := m.itemService.SleepItem(item); err != nil {
						return fmt.Errorf("sleep item: %w", err)
					} else {
						return itemUpdatedMsg{item: item}
					}
				}
			case "awake":
				return m, func() tea.Msg {
					if item, err := m.itemService.AwakeItem(item); err != nil {
						return fmt.Errorf("awake item: %w", err)
					} else {
						return itemUpdatedMsg{item: item}
					}
				}
			case "close":
				return m, func() tea.Msg {
					if err := m.itemService.CloseItem(item); err != nil {
						return fmt.Errorf("close item: %w", err)
					}
					return itemClosedMsg{}
				}
			case "delete":
				return m, func() tea.Msg {
					if err := m.itemService.DeleteItem(item); err != nil {
						return fmt.Errorf("delete item: %w", err)
					}
					return itemDeletedMsg{}
				}
			case "item_type":
				m.currentScreen = screenItemType
			case "item_nodes":
				m.currentScreen = screenItemNodes
			case "item_notes":
				m.currentScreen = screenItemNotes
			case "item_links":
				m.currentScreen = screenItemLinks
			case "item_vms":
				m.currentScreen = screenItemVMs
			}
		case "s":
			if item.IsSleep() {
				return m, func() tea.Msg {
					if item, err := m.itemService.AwakeItem(item); err != nil {
						return fmt.Errorf("awake: %w", err)
					} else {
						return itemUpdatedMsg{item: item}
					}
				}
			} else {
				return m, func() tea.Msg {
					if item, err := m.itemService.SleepItem(item); err != nil {
						return fmt.Errorf("sleep: %w", err)
					} else {
						return itemUpdatedMsg{item: item}
					}
				}
			}
		}
	}

	return m, nil
}

func (m *TeaModel) viewEditItem() string {
	var state string

	item, found := m.getSelectedItem()
	if !found {
		return "item not found\n"
	}

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
	item, found := m.getSelectedItem()
	if !found {
		return
	}

	m.editItemMenu.ResetMenu()

	m.editItemMenu.AddGroup("exit")

	if !item.IsClosed() {
		m.editItemMenu.AddGroup("item_type")
	}

	m.editItemMenu.AddGroup("item_nodes")
	m.editItemMenu.AddGroup("item_vms")
	m.editItemMenu.AddGroup("item_notes")
	m.editItemMenu.AddGroup("item_links")
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

	if g, _ := m.editItemMenu.GetGroup(); g == "" {
		m.editItemMenu.JumpToGroup("exit")
	}
}

func (m *TeaModel) getSelectedItem() (_ model.Item, found bool) {
	for i := range m.items {
		if m.items[i].ID == m.selectedItemID {
			return m.items[i], true
		}
	}
	return model.Item{}, false
}
