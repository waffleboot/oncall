package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
	"github.com/waffleboot/oncall/pkg/tea/button"
	"github.com/waffleboot/oncall/pkg/tea/menu"
	"github.com/waffleboot/oncall/pkg/tea/tabs"
	"go.uber.org/zap"
)

const (
	screenItems    screen = "items"
	screenItem     screen = "item"
	screenItemType screen = "item_type"
	screenLinks    screen = "links"
	screenLink     screen = "link"
	screenNodes    screen = "nodes"
	screenNode     screen = "node"
	screenNewNodes screen = "new_nodes"
	screenNotes    screen = "notes"
	screenNote     screen = "note"
	screenVMs      screen = "vms"
	screenVM       screen = "vm"
	screenTitle    screen = "title"
	screenUsers    screen = "users"
)

type (
	screen   string
	TeaModel struct {
		userService              port.UserService
		itemService              port.ItemService
		journalService           port.JournalService
		currentScreen            screen
		users                    []model.User
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
		menuAllItems             *menu.Model
		menuItem                 *menu.Model
		menuItemType             *menu.Model
		linksMenu                *menu.Model
		menuVMs                  *menu.Model
		menuNotes                *menu.Model
		menuNodes                *menu.Model
		menuUsers                *menu.Model
		tabs                     tabs.Model
		textinputLinkAddress     textinput.Model
		textinputLinkDescription textarea.Model
		textinputItemTitle       textinput.Model
		textinputItemDescription textarea.Model
		textinputVmName          textinput.Model
		textinputVmNode          textinput.Model
		textinputVmDescription   textarea.Model
		textinputNode            textinput.Model
		textinputNodes           textarea.Model
		textinputNote            textarea.Model
		submitVM                 button.Model
		submitTitle              button.Model
		submitLink               button.Model
		submitNote               button.Model
		submitNodes              button.Model
		submitAsPublicLink       button.Model
		submitAsPrivateLink      button.Model
		printJournal             bool
		log                      *zap.Logger
		err                      error
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

func NewTeaModel(
	userService port.UserService,
	itemService port.ItemService,
	journalService port.JournalService,
	users []model.User,
	log *zap.Logger) *TeaModel {
	return &TeaModel{
		userService:    userService,
		itemService:    itemService,
		journalService: journalService,
		users:          users,
		log:            log,
	}
}

func (m *TeaModel) Init() tea.Cmd {
	m.currentScreen = screenItems
	m.menuAllItems = menu.New(func(group string, pos int) string {
		switch {
		case group == "exit":
			return "Exit"
		case group == "new":
			return "Новое обращение"
		case group == "close_journal":
			return "Закрыть журнал"
		case group == "print_journal":
			return "Распечатать журнал"
		case group == "next":
			if user := m.userService.GetUser(); user != nil {
				return fmt.Sprintf("Следующий дежурный - %s ...", user.Name)
			} else {
				return "Set next ..."
			}
		case group == "items":
			marker := " "
			item := m.items[pos]
			switch {
			case item.IsSleep():
				marker = "?"
			case item.IsClosed():
				marker = "x"
			}
			return fmt.Sprintf("%s #%d - %s - %s", marker, item.Num, item.Type, item.MenuItem())
		}
		return group
	})
	m.resetItems(nil)
	m.menuItem = menu.New(func(group string, pos int) string {
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
		case group == string(screenTitle):
			return "Имя и описание ..."
		case group == string(screenItemType):
			return fmt.Sprintf("Тип обращения: (%s) ...", m.selectedItem.Type)
		case group == string(screenNotes):
			return "Заметки ..."
		case group == string(screenLinks):
			return "Ссылки ..."
		case group == string(screenNodes):
			return "Хосты, узлы ..."
		case group == string(screenVMs):
			return "ВМ-ки ..."
		}
		return group
	})
	m.menuItemType = menu.New(func(group string, pos int) string {
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
	m.linksMenu = menu.New(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить ссылку ..."
		case "links":
			return m.links[pos].MenuItem()
		default:
			return group
		}
	})
	m.menuVMs = menu.New(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить ВМ ..."
		case "vms":
			return m.vms[pos].MenuItem()
		default:
			return group
		}
	})
	m.menuNodes = menu.New(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить узлы ..."
		case "nodes":
			node := m.nodes[pos]

			var s strings.Builder
			s.WriteString(fmt.Sprintf("#%d - ", node.ID))
			s.WriteString(node.MenuItem())
			return s.String()
		default:
			return group
		}
	})
	m.menuNotes = menu.New(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "new":
			return "Добавить заметку ..."
		case "notes":
			return m.notes[pos].MenuItem()
		default:
			return group
		}
	})
	m.menuUsers = menu.New(func(group string, pos int) string {
		switch group {
		case "exit":
			return "Exit"
		case "users":
			return m.users[pos].MenuItem()
		}
		return ""
	})
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+p":
			m.printJournal = true
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, tea.Quit
	}

	switch m.currentScreen {
	case screenItems:
		return m.updateItems(msg)
	case screenItem:
		return m.updateItem(msg)
	case screenItemType:
		return m.updateItemType(msg)
	case screenNodes:
		return m.updateNodes(msg)
	case screenNode:
		return m.updateNode(msg)
	case screenNewNodes:
		return m.updateNewNodes(msg)
	case screenNotes:
		return m.updateNotes(msg)
	case screenNote:
		return m.updateNote(msg)
	case screenLinks:
		return m.updateLinks(msg)
	case screenLink:
		return m.updateLink(msg)
	case screenVMs:
		return m.updateVMs(msg)
	case screenVM:
		return m.updateVM(msg)
	case screenTitle:
		return m.updateItemTitle(msg)
	case screenUsers:
		return m.updateUsers(msg)
	}

	return m, func() tea.Msg { return fmt.Errorf("screen not found: %s", m.currentScreen) }
}

func (m *TeaModel) View() string {
	switch m.currentScreen {
	case screenItems:
		return m.viewItems()
	case screenItem:
		return m.viewItem()
	case screenItemType:
		return m.viewItemType()
	case screenNodes:
		return m.viewNodes()
	case screenNode:
		return m.viewNode()
	case screenNewNodes:
		return m.viewNewNodes()
	case screenNotes:
		return m.viewNotes()
	case screenNote:
		return m.viewNote()
	case screenLinks:
		return m.viewLinks()
	case screenLink:
		return m.viewLink()
	case screenVMs:
		return m.viewVMs()
	case screenVM:
		return m.viewVM()
	case screenTitle:
		return m.viewTitle()
	case screenUsers:
		return m.viewUsers()
	}
	return string(m.currentScreen)
}

func (m *TeaModel) getItems() tea.Msg {
	return m.itemService.GetItems()
}

func (m *TeaModel) getItem() tea.Msg {
	item, err := m.itemService.GetItem(m.selectedItem.ID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	return item
}

func (m *TeaModel) PrintJournal() bool {
	return m.printJournal
}

func (m *TeaModel) Err() error {
	return m.err
}

func (m *TeaModel) exitScreen() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg { return "exit" }
}

func (m *TeaModel) runAndExitScreen(f func() error) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		if err := f(); err != nil {
			return err
		}
		return "exit"
	}
}
