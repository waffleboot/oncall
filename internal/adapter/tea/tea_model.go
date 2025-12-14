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
		selectedItem  int
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
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
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

func (m TeaModel) getItems() tea.Msg {
	items, err := m.config.ItemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}

func (m *TeaModel) itemLabel(item model.Item) string {
	return fmt.Sprintf("  #%d - %s", item.ID, item.Type)
}
