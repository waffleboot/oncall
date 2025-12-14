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
	screen   string
	TeaModel struct {
		itemService   port.ItemService
		currentScreen screen
		items         []model.Item
		selectedItem  int
		allItemsMenu  *Menu
		editItemMenu  *Menu
	}
	newItemCreatedMsg struct {
		newItem model.Item
	}
	itemUpdatedMsg struct{}
)

func NewTeaModel(itemService port.ItemService) *TeaModel {
	return &TeaModel{itemService: itemService}
}

func (m *TeaModel) Init() tea.Cmd {
	m.currentScreen = screenAllItems
	m.allItemsMenu = NewMenu(func(group string, pos int) string {
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
	m.editItemMenu = NewMenu(func(group string, pos int) string {
		switch {
		case group == "exit":
			return "Exit"
		case group == "sleep":
			return "В ожидание"
		case group == "awake":
			return "Из ожидания"
		case group == "close":
			return "Закрыть"
		case group == "delete":
			return "Удалить"
		case group == "edit_type":
			return fmt.Sprintf("Тип обращения: (%s)...", m.items[m.selectedItem].Type)
		case group == "notes":
			return "Заметки..."
		case group == "links":
			return "Ссылки..."
		case group == "nodes":
			return "Хосты, узлы..."
		case group == "vms":
			return "ВМ-ки..."
		}
		return ""
	})
	return func() tea.Msg {
		items, err := m.itemService.GetItems()
		if err != nil {
			return fmt.Errorf("get items: %w", err)
		}
		return items
	}
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
		m.resetAllItemsMenu()
	case newItemCreatedMsg:
		m.items = append(m.items, msg.newItem)
		m.resetAllItemsMenu()
		m.selectedItem = len(m.items) - 1
		m.allItemsMenu.JumpToPos("items", m.selectedItem)
	case itemUpdatedMsg:
		m.resetAllItemsMenu()
		m.resetEditItemMenu()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentScreen {
	case screenAllItems:
		return m.updateAllItems(msg)
	case screenEditItem:
		return m.updateEditItem(msg)
	}
	return m, nil
}

func (m *TeaModel) View() string {
	switch m.currentScreen {
	case screenAllItems:
		return m.viewAllItems()
	case screenEditItem:
		return m.viewEditItem()
	}
	return ""
}
