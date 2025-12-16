package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	screenAllItems screen = "all_items"
	screenItem     screen = "edit_item"
	screenItemType screen = "item_type"
	screenLinks    screen = "links"
	screenLink     screen = "link"
	screenNodes    screen = "nodes"
	screenNode     screen = "node"
	screenNotes    screen = "notes"
	screenNote     screen = "note"
	screenVMs      screen = "vms"
	screenVM       screen = "vm"
	screenTitle    screen = "title"
)

type (
	screen   string
	TeaModel struct {
		itemService              port.ItemService
		journalService           port.JournalService
		currentScreen            screen
		items                    []model.Item
		vms                      []model.VM
		links                    []model.Link
		nodes                    []model.Node
		notes                    []model.Note
		selectedItem             model.Item
		selectedLink             model.Link
		selectedVM               model.VM
		selectedNode             model.Node
		selectedNote             model.Note
		menuAllItems             *Menu
		menuEditItem             *Menu
		menuItemType             *Menu
		linksMenu                *Menu
		menuVMs                  *Menu
		notesVMs                 *Menu
		nodesVMs                 *Menu
		textinputLinkAddress     textinput.Model
		textinputLinkDescription textarea.Model
		textinputItemTitle       textinput.Model
		textinputItemDescription textarea.Model
		textinputVmName          textinput.Model
		textinputVmNode          textinput.Model
		textinputVmDescription   textarea.Model
		textinputNodeName        textinput.Model
		textinputNote            textarea.Model
		textInput                string
		printJournal             bool
	}
	itemCreatedMsg struct {
		item model.Item
	}
	itemUpdatedMsg     struct{}
	itemClosedMsg      struct{}
	itemDeletedMsg     struct{}
	itemLinkCreatedMsg struct {
		link model.Link
	}
)

func NewTeaModel(itemService port.ItemService, journalService port.JournalService) *TeaModel {
	return &TeaModel{itemService: itemService, journalService: journalService}
}

func (m *TeaModel) Init() tea.Cmd {
	m.currentScreen = screenAllItems
	m.menuAllItems = NewMenu(func(group string, pos int) string {
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
			marker := " "
			item := m.items[pos]
			switch {
			case item.IsSleep():
				marker = "?"
			case item.IsClosed():
				marker = "x"
			}
			return fmt.Sprintf("%s #%d - %s - %s", marker, item.Num, item.Type, item.TitleForView())
		}
		return group
	})
	m.resetAllItems(nil)
	m.menuEditItem = NewMenu(func(group string, pos int) string {
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
		case group == "item_title":
			return "Имя и описание ..."
		case group == "item_type":
			return fmt.Sprintf("Тип обращения: (%s) ...", m.selectedItem.Type)
		case group == "item_notes":
			return "Заметки ..."
		case group == "item_links":
			return "Ссылки ..."
		case group == "item_nodes":
			return "Хосты, узлы ..."
		case group == "item_vms":
			return "ВМ-ки ..."
		}
		return group
	})
	m.menuItemType = NewMenu(func(group string, pos int) string {
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
	m.menuItemType.AddGroup(string(model.ItemTypeAsk))
	m.menuItemType.AddGroup(string(model.ItemTypeInc))
	m.menuItemType.AddGroup(string(model.ItemTypeAlert))
	m.menuItemType.AddGroup(string(model.ItemTypeAdhoc))
	m.linksMenu = NewMenu(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить ссылку ..."
		case "links":
			link := m.links[pos]

			var s strings.Builder
			s.WriteString(fmt.Sprintf("#%d - ", link.ID))
			if link.Address == "" {
				s.WriteString("empty")
			} else {
				s.WriteString(link.Address)
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
	m.menuVMs = NewMenu(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить ВМ ..."
		case "vms":
			vm := m.vms[pos]

			var s strings.Builder
			s.WriteString(fmt.Sprintf("#%d - ", vm.ID))
			if vm.Name == "" {
				s.WriteString("empty")
			} else {
				s.WriteString(vm.Name)
			}
			return s.String()
		default:
			return group
		}
	})
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
	case screenItem:
		return m.updateEditItem(msg)
	case screenItemType:
		return m.updateItemType(msg)
	case screenNodes:
		return m.updateItemNodes(msg)
	case screenNotes:
		return m.updateItemNotes(msg)
	case screenLinks:
		return m.updateItemLinks(msg)
	case screenLink:
		return m.updateItemLink(msg)
	case screenVMs:
		return m.updateVMs(msg)
	case screenVM:
		return m.updateVM(msg)
	case screenTitle:
		return m.updateItemTitle(msg)
	}
	return m, nil
}

func (m *TeaModel) View() string {
	switch m.currentScreen {
	case screenAllItems:
		return m.viewAllItems()
	case screenItem:
		return m.viewEditItem()
	case screenItemType:
		return m.viewItemType()
	case screenNodes:
		return m.viewItemNodes()
	case screenNotes:
		return m.viewItemNotes()
	case screenLinks:
		return m.viewItemLinks()
	case screenLink:
		return m.viewItemLink()
	case screenVMs:
		return m.viewVMs()
	case screenVM:
		return m.viewVM()
	case screenTitle:
		return m.viewItemTitle()
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

func (m *TeaModel) PrintJournal() bool {
	return m.printJournal
}
