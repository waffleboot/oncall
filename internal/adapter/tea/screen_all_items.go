package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

type allItemsModel struct {
	menu *Menu
}

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsModel.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case []model.Item:
		m.allItemsModel.menu.ResetMenu()
		m.allItemsModel.menu.AddGroup("exit")
		m.allItemsModel.menu.AddGroup("new")
		m.allItemsModel.menu.AddGroup("close_journal")
		m.allItemsModel.menu.AddGroup("print_journal")
		m.allItemsModel.menu.AddGroupWithItems("items", len(m.items))
		m.allItemsModel.menu.AdjustCursor()
		m.allItemsModel.menu.JumpToItem("items", func(pos int) (found bool) {
			return m.items[pos].ID == m.selectedItemID
		})
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
					m.selectedItemID = item.ID
					return m.getItems()
				}
			case "close_journal":
			case "print_journal":
			case "items":
				m.selectedItemID = m.items[p].ID
				m.screenPush(screenEditItem)
			}
		}
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.allItemsModel.menu.GenerateMenu()
}
