package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	screenAllItems screen = "all_items"
)

type (
	screen         string
	TeaModelConfig struct {
		ItemService port.ItemService
	}
	TeaModel struct {
		config        TeaModelConfig
		screens       []screen
		items         []model.Item
		allItemsModel AllItemsModel
	}
	AllItemsModel struct {
		menu Menu
	}
)

func NewTeaModel(config TeaModelConfig) TeaModel {
	return TeaModel{config: config}
}

func (m TeaModel) Init() tea.Cmd {
	return m.getItems
}

func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.topScreen() {
	case screenAllItems:
		return m.updateAllItems(msg)
	}
	return m, nil
}

func (m TeaModel) View() string {
	switch m.topScreen() {
	case screenAllItems:
		return m.viewAllItems()
	}
	return ""
}

func (m TeaModel) viewAllItems() string {
	return m.allItemsModel.menu.GenerateMenu()
}

func (m TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsModel.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
		m.allItemsModel.menu = NewMenu(func(group string, pos int) string {
			switch {
			case group == "new":
				return "Новое обращение"
			case group == "items":
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
		m.allItemsModel.menu.AddGroupWithItems("items", len(m.items))
		return m, nil
	}
	return m, nil
}

func (m TeaModel) topScreen() screen {
	return m.screens[len(m.screens)-1]
}

func (m TeaModel) pushScreen(screen screen) []screen {
	return append(m.screens, screen)
}

func (m TeaModel) popScreen() []screen {
	if len(m.screens) == 1 {
		return m.screens
	}
	return m.screens[:len(m.screens)-1]
}

func (m TeaModel) getItems() tea.Msg {
	items, err := m.config.ItemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}
