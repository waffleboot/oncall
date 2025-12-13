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
	next           tea.Model
}

func NewModelCloseJournal(controller *Controller, journalService port.JournalService, next tea.Model) *ModelCloseJournal {
	m := &ModelCloseJournal{controller: controller, journalService: journalService, next: next}
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
			return m.next, nil
		case "enter", " ":
			switch g, _ := m.menu.GetGroup(); g {
			case closeJournalYes:
				return m.next, func() tea.Msg {
					if err := m.journalService.CloseJournal(); err != nil {
						return fmt.Errorf("close journal: %w", err)
					}
					return "journal closed"
				}
			case closeJournalNo:
				return m.next, nil
			}
		}
	}
	return m, nil
}

func (m *ModelCloseJournal) View() string {
	var s strings.Builder
	s.WriteString("  Закрыть журнал?\n\n")
	s.WriteString(m.menu.GenerateMenu())
	return s.String()
}
