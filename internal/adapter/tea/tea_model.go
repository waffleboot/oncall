package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	screenItemLink  screen = "item_link"
)

type (
	screen   string
	TeaModel struct {
		itemService       port.ItemService
		journalService    port.JournalService
		currentScreen     screen
		items             []model.Item
		links             []model.ItemLink
		selectedItem      model.Item
		selectedLink      model.ItemLink
		allItemsMenu      *Menu
		editItemMenu      *Menu
		editItemTypeMenu  *Menu
		editItemLinksMenu *Menu
		linkInput         textinput.Model
	}
	itemCreatedMsg struct {
		item model.Item
	}
	itemUpdatedMsg     struct{}
	itemClosedMsg      struct{}
	itemDeletedMsg     struct{}
	itemLinkCreatedMsg struct {
		link model.ItemLink
	}
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
			item := m.items[pos]
			switch {
			case item.IsSleep():
				return fmt.Sprintf("? #%d - %s", item.Num, item.Type)
			case item.IsClosed():
				return fmt.Sprintf("x #%d - %s", item.Num, item.Type)
			default:
				return fmt.Sprintf("  #%d - %s", item.Num, item.Type)
			}
		}
		return group
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
		return group
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
		return group
	})
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeInc))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAsk))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAlert))
	m.editItemTypeMenu.AddGroup(string(model.ItemTypeAdhoc))
	m.editItemLinksMenu = NewMenu(func(group string, pos int) string {
		switch group {
		case "new":
			return "Добавить ссылку..."
		case "links":
			link := m.links[pos]

			var s strings.Builder
			s.WriteString(fmt.Sprintf("#%d - ", link.ID))
			if link.Link == "" {
				s.WriteString("empty")
			} else {
				s.WriteString(link.Link)
			}
			if link.Public {
				s.WriteString(" - public")
			} else {
				s.WriteString(" - private")
			}
			return s.String()
		default:
			return group
		}
	})
	m.linkInput = textinput.New()
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	case screenItemLink:
		return m.updateItemLink(msg)
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
	case screenItemLink:
		return m.viewItemLink()
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

func (m *TeaModel) getItem() tea.Msg {
	item, err := m.itemService.GetItem(m.selectedItem.ID)
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return item
}
