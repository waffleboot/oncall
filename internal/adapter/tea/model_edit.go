package tea

import (
	"fmt"
	"strings"

	"github.com/waffleboot/oncall/internal/model"

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

type ModelEdit struct {
	controller  *Controller
	itemService port.ItemService
	itemID      int
	item        model.Item
	menu        *Menu
}

func NewModelEdit(controller *Controller, itemService port.ItemService, itemID int) *ModelEdit {
	m := &ModelEdit{controller: controller, itemService: itemService, itemID: itemID}
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
	return m
}

func (m *ModelEdit) Init() tea.Cmd {
	return m.getItem
}

func (m *ModelEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "e":
			next := m.controller.modelStart()
			return next, next.Init()
		case "enter", " ":
			switch g, _ := m.menu.GetGroup(); g {
			case editExit:
				next := m.controller.modelStart()
				return next, next.Init()
			case editSleep:
				return m, func() tea.Msg {
					if err := m.itemService.SleepItem(m.item); err != nil {
						return fmt.Errorf("sleep item: %w", err)
					}
					return m.getItem()
				}
			case editAwake:
				return m, func() tea.Msg {
					if err := m.itemService.AwakeItem(m.item); err != nil {
						return fmt.Errorf("awake item: %w", err)
					}
					return m.getItem()
				}
			case editClose:
				return m, func() tea.Msg {
					if err := m.itemService.CloseItem(m.item); err != nil {
						return fmt.Errorf("close item: %w", err)
					}
					return "closed"
				}
			case editDelete:
				return m, func() tea.Msg {
					if err := m.itemService.DeleteItem(m.item); err != nil {
						return fmt.Errorf("delete item: %w", err)
					}
					return "deleted"
				}
			case editType:
				next := m.controller.modelItemType(m.item)
				return next, next.Init()
			}
		}
	case error:
		return m.controller.modelError(msg.Error(), m), nil
	case model.Item:
		m.item = msg
		m.resetMenu()
		return m, nil
	case string:
		if msg == "closed" || msg == "deleted" {
			next := m.controller.modelStart()
			return next, next.Init()
		}
	}
	return m, nil
}

func (m *ModelEdit) View() string {
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

func (m *ModelEdit) resetMenu() {
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

func (m *ModelEdit) getItem() tea.Msg {
	item, err := m.itemService.GetItem(m.itemID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	return item
}
