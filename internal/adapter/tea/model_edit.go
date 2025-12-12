package tea

import (
	"fmt"
	"github.com/waffleboot/oncall/internal/model"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	editExit   = "exit"
	editSleep  = "sleep"
	editAwake  = "awake"
	editClose  = "close"
	editDelete = "delete"
	editType   = "type"
	editNotes  = "notes"
	editLinks  = "links"
	editNodes  = "nodes"
	editVMs    = "vms"
)

type EditModel struct {
	controller  *Controller
	itemService port.ItemService
	item        model.Item
	prev        Prev
	menu        *Menu
}

func NewEditModel(controller *Controller, itemService port.ItemService, item model.Item, prev Prev) *EditModel {
	m := &EditModel{controller: controller, itemService: itemService, prev: prev, item: item}

	m.menu = NewMenu(func(group string, pos int) string {
		switch {
		case group == editExit:
			return "Exit"
		case group == editSleep:
			return "В ожидание"
		case group == editAwake:
			return "Из ожидания"
		case group == editClose:
			return "Закрыть"
		case group == editDelete:
			return "Удалить"
		case group == editType:
			return fmt.Sprintf("Тип обращения: (%s)...", m.item.Type)
		case group == editNotes:
			return "Заметки..."
		case group == editLinks:
			return "Ссылки..."
		case group == editNodes:
			return "Хосты, узлы..."
		case group == editVMs:
			return "ВМ-ки..."
		}
		return ""
	})

	m.resetMenu(item)

	return m
}

func (m *EditModel) Init() tea.Cmd {
	return nil
}

func (m *EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m.prev()
		case "enter", " ":
			switch g, _ := m.menu.GetGroup(); g {
			case editExit:
				return m.prev()
			case editSleep:
				if err := m.itemService.SleepItem(m.item); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.resetItem()
			case editAwake:
				if err := m.itemService.AwakeItem(m.item); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.resetItem()
			case editClose:
				if err := m.itemService.CloseItem(m.item); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.prev()
			case editDelete:
				if err := m.itemService.DeleteItem(m.item); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.prev()
			case editType:
				next := m.controller.itemTypeModel(m.item, m.resetItem)
				return next, next.Init()
			}
		}
	}
	return m, nil
}

func (m *EditModel) View() string {
	var state string

	switch {
	case m.item.IsSleep():
		state = " в ожидании"
	case m.item.IsClosed():
		switch m.item.Type {
		case model.ItemTypeAsk:
			state = " закрыто"
		default:
			state = " закрыт"
		}
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("  #%d %s%s\n\n", m.item.ID, m.item.Type, state))
	s.WriteString(m.menu.GenerateMenu())

	return s.String()
}

func (m *EditModel) resetItem() (tea.Model, tea.Cmd) {
	item, err := m.itemService.GetItem(m.item.ID)
	if err != nil {
		return m.controller.errorModel(err.Error(), m.prev), nil
	}

	m.resetMenu(item)

	return m, nil
}

func (m *EditModel) resetMenu(item model.Item) {
	m.item = item

	m.menu.ResetMenu()

	m.menu.AddGroup(editExit)

	if !m.item.IsClosed() {
		m.menu.AddGroup(editType)
	}

	m.menu.AddGroup(editNodes)
	m.menu.AddGroup(editVMs)
	m.menu.AddGroup(editNotes)
	m.menu.AddGroup(editLinks)
	m.menu.AddDelimiter()

	if m.item.IsActive() {
		m.menu.AddGroup(editSleep)
	}

	if m.item.IsSleep() {
		m.menu.AddGroup(editAwake)
	}

	if !m.item.IsClosed() {
		m.menu.AddGroup(editClose)
	}

	m.menu.AddDelimiter()
	m.menu.AddGroup(editDelete)
}
