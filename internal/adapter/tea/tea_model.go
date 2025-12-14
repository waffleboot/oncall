package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	screenAllItems screen = "all_items"
	screenEditItem screen = "edit_item"
)

type (
	TeaModelConfig struct {
		ItemService port.ItemService
	}
	TeaModel struct {
		config        TeaModelConfig
		screens       []screen
		items         []model.Item
		selectedItem  model.Item
		allItemsModel allItemsModel
		editItemModel editItemModel
	}
	editItemModel struct{}
)

func NewTeaModel(config TeaModelConfig) *TeaModel {
	return &TeaModel{config: config}
}

func (m *TeaModel) Init() tea.Cmd {
	m.screenPush(screenAllItems)
	m.allItemsModel.menu = NewMenu(func(group string, pos int) string {
		switch {
		case group == "exit":
			return "Exit"
		case group == "new":
			return "Новое обращение"
		case group == "close_journal":
			return "Закрыть журнал"
		case group == "print_journal":
			return "Распечатать журнал"
		case group == "items":
			return m.itemLabel(m.items[pos])
		}
		return ""
	})
	return m.getInitialItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
		m.allItemsModel.menu.AddGroup("exit")
		m.allItemsModel.menu.AddGroup("new")
		m.allItemsModel.menu.AddGroup("close_journal")
		m.allItemsModel.menu.AddGroup("print_journal")
		m.allItemsModel.menu.AddGroupWithItems("items", len(m.items))
	case newItemCreateMsg:
		m.items = msg.items
		m.selectedItem = msg.newItem
		m.allItemsModel.menu.ResetMenu()
		m.allItemsModel.menu.AddGroup("exit")
		m.allItemsModel.menu.AddGroup("new")
		m.allItemsModel.menu.AddGroup("close_journal")
		m.allItemsModel.menu.AddGroup("print_journal")
		m.allItemsModel.menu.AddGroupWithItems("items", len(m.items))
		m.allItemsModel.menu.JumpToItem("items", func(pos int) (found bool) {
			return m.items[pos].ID == m.selectedItem.ID
		})
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentScreen() {
	case screenAllItems:
		return m.updateAllItems(msg)
	case screenEditItem:
		return m.updateEditItem(msg)
	}
	return m, nil
}

func (m *TeaModel) View() string {
	switch m.currentScreen() {
	case screenAllItems:
		return m.viewAllItems()
	case screenEditItem:
		return m.viewEditItem()
	}
	return ""
}

func (m TeaModel) getInitialItems() tea.Msg {
	items, err := m.config.ItemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}

func (m *TeaModel) itemLabel(item model.Item) string {
	return fmt.Sprintf("  #%d - %s", item.ID, item.Type)
}
