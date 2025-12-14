package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

type newItemCreatedMsg struct {
	newItem model.Item
}

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		case "enter", " ":
			switch g, p := m.allItemsMenu.GetGroup(); g {
			case "exit":
				return m, tea.Quit
			case "new":
				return m, func() tea.Msg {
					item := m.itemService.CreateItem()
					if err := m.itemService.UpdateItem(item); err != nil {
						return fmt.Errorf("create item: %w", err)
					}

					return newItemCreatedMsg{newItem: item}
				}
			case "close_journal":
			case "print_journal":
			case "items":
				m.selectedItem = p
				m.currentScreen = screenEditItem
				m.resetEditItemMenu()
			}
		}
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.allItemsMenu.GenerateMenu()
}
