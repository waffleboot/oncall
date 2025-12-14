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
	m := &TeaModel{config: config}
	m.screenPush(screenAllItems)
	return m
}

func (m *TeaModel) Init() tea.Cmd {
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
		m.allItemsModel.menu = NewMenu(func(group string, pos int) string {
			switch {
			case group == startNew:
				return "Новое обращение"
			case group == startItems:
				return m.itemLabel(m.items[pos])
			case group == startClose:
				return "Закрыть журнал"
			case group == startPrint:
				return "Распечатать журнал"
			case group == startExit:
				return "Exit"
			}
			return ""
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

func (m TeaModel) getItems() tea.Msg {
	items, err := m.config.ItemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}
