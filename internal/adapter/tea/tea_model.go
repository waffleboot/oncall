package tea

import (
	"cmp"
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	screenAllItems  screen = "all_items"
	screenEditItem  screen = "edit_item"
	screenItemType  screen = "item_type"
	screenItemLinks screen = "item_links"
	screenItemNodes screen = "item_nodes"
	screenItemNotes screen = "item_notes"
	screenItemVMs   screen = "item_vms"
)

type (
	screen   string
	TeaModel struct {
		itemService      port.ItemService
		journalService   port.JournalService
		currentScreen    screen
		items            []model.Item
		selectedItem     model.Item
		allItemsMenu     *Menu
		editItemMenu     *Menu
		editItemTypeMenu *Menu
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
		case group == "item_type":
			return fmt.Sprintf("Тип обращения: (%s)...", m.selectedItem.Type)
		case group == "item_notes":
			return "Заметки..."
		case group == "item_links":
			return "Ссылки..."
		case group == "item_nodes":
			return "Хосты, узлы..."
		case group == "item_vms":
			return "ВМ-ки..."
		}
		return ""
	})
	m.editItemTypeMenu = NewMenu(func(group string, pos int) string {
		switch model.ItemType(group) {
		case model.ItemTypeInc:
			return "Инцидент"
		case model.ItemTypeAdhoc:
			return "Adhoc"
		case model.ItemTypeAsk:
			return "Обращение"
		case model.ItemTypeAlert:
			return "Alert"
		}
		return ""
	})
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeInc))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAsk))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAlert))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAdhoc))
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg

		slices.SortFunc(m.items, func(a, b model.Item) int {
			if c := a.Type.Compare(b.Type); c != 0 {
				return c
			}
			return cmp.Compare(a.ID, b.ID)
		})

		m.resetAllItemsMenu()
	case itemCreatedMsg:
		m.selectedItem = msg.item
		return m, m.getItems
	case itemUpdatedMsg:
		m.selectedItem = msg.item
		m.resetEditItemMenu()
		if m.currentScreen == screenItemType {
			m.currentScreen = screenEditItem
		}
		return m, m.getItems
	case itemClosedMsg:
		m.currentScreen = screenAllItems
		return m, m.getItems
	case itemDeletedMsg:
		m.currentScreen = screenAllItems
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
	case screenItemType:
		return m.updateItemType(msg)
	case screenItemNodes:
		return m.updateItemNodes(msg)
	case screenItemNotes:
		return m.updateItemNotes(msg)
	case screenItemLinks:
		return m.updateItemLinks(msg)
	case screenItemVMs:
		return m.updateItemVMs(msg)
	}
	return m, nil
}

func (m *TeaModel) View() string {
	switch m.currentScreen {
	case screenAllItems:
		return m.viewAllItems()
	case screenEditItem:
		return m.viewEditItem()
	case screenItemType:
		return m.viewItemType()
	case screenItemNodes:
		return m.viewItemNodes()
	case screenItemNotes:
		return m.viewItemNotes()
	case screenItemLinks:
		return m.viewItemLinks()
	case screenItemVMs:
		return m.viewItemVMs()
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
