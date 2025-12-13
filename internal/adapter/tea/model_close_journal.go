package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	closeJournalYes = "yes"
	closeJournalNo  = "no"
)

type ModelCloseJournal struct {
	controller     *Controller
	journalService port.JournalService
	menu           *Menu
}

func NewModelCloseJournal(controller *Controller, journalService port.JournalService) *ModelCloseJournal {
	m := &ModelCloseJournal{controller: controller, journalService: journalService}
	m.menu = NewMenu(func(group string, pos int) string {
		switch group {
		case closeJournalYes:
			return "Yes"
		case closeJournalNo:
			return "No"
		default:
			return ""
		}
	})
	m.menu.AddGroup(closeJournalNo)
	m.menu.AddGroup(closeJournalYes)
	return m
}

func (m *ModelCloseJournal) Init() tea.Cmd {
	return nil
}

func (m *ModelCloseJournal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			next := m.controller.modelStart()
			return next, next.Init()
		case "enter", " ":
			switch g, _ := m.menu.GetGroup(); g {
			case closeJournalYes:
				return m, func() tea.Msg {
					if err := m.journalService.CloseJournal(); err != nil {
						return fmt.Errorf("close journal: %w", err)
					}
					return "journal closed"
				}
			case closeJournalNo:
				next := m.controller.modelStart()
				return next, next.Init()
			}
		}
	case string:
		if msg == "journal closed" {
			next := m.controller.modelStart()
			return next, next.Init()
		}
	case error:
		return m.controller.modelError(msg.Error(), m), nil
	}
	return m, nil
}

func (m *ModelCloseJournal) View() string {
	var s strings.Builder
	s.WriteString("  Закрыть журнал?\n\n")
	s.WriteString(m.menu.GenerateMenu())
	return s.String()
}
