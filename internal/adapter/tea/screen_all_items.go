package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

type (
	allItemsModel struct {
		menu *Menu
	}
	newItemCreatedMsg struct {
		items   []model.Item
		newItem model.Item
	}
)

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsModel.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		case "enter", " ":
			switch g, p := m.allItemsModel.menu.GetGroup(); g {
			case "exit":
				return m, tea.Quit
			case "new":
				return m, func() tea.Msg {
					item := m.config.ItemService.CreateItem()
					if err := m.config.ItemService.UpdateItem(item); err != nil {
						return fmt.Errorf("create item: %w", err)
					}

					items, err := m.config.ItemService.GetItems()
					if err != nil {
						return fmt.Errorf("get items: %w", err)
					}

					return newItemCreatedMsg{items: items, newItem: item}
				}
			case "close_journal":
			case "print_journal":
			case "items":
				m.selectedItem = m.items[p]
				m.screenPush(screenEditItem)
			}
		}
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.allItemsModel.menu.GenerateMenu()
}
