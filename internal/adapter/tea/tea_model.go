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
		itemService    port.ItemService
		journalService port.JournalService
		currentScreen  screen
		items          []model.Item
		selectedItem   int
		allItemsMenu   *Menu
		editItemMenu   *Menu
	}
	itemCreatedMsg struct {
		item model.Item
	}
	itemUpdatedMsg struct {
		item model.Item
	}
	itemClosedMsg  struct{}
	itemDeletedMsg struct{}
)

func NewTeaModel(itemService port.ItemService, journalService port.JournalService) *TeaModel {
	return &TeaModel{itemService: itemService, journalService: journalService}
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
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
		m.resetAllItemsMenu()
	case itemCreatedMsg:
		m.items = append(m.items, msg.item)
		m.resetAllItemsMenu()
		m.allItemsMenu.JumpToPos("items", len(m.items)-1)
	case itemUpdatedMsg:
		m.items[m.selectedItem] = msg.item
		m.resetEditItemMenu()
		return m, nil
	case itemClosedMsg:
		m.currentScreen = screenAllItems
		return m, m.getItems
	case itemDeletedMsg:
		m.currentScreen = screenAllItems
		if len(m.items) == 0 {
			m.allItemsMenu.JumpToGroup("new")
		}
		return m, m.getItems
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

func (m *TeaModel) getItems() tea.Msg {
	items, err := m.itemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}
