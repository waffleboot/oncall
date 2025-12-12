package tea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/port"
	"strings"
)

const (
	closeJournalYes = "yes"
	closeJournalNo  = "no"
)

type CloseJournalModel struct {
	controller     *Controller
	journalService port.JournalService
	prev           Prev
	menu           *Menu
}

func NewCloseJournalModel(controller *Controller, journalService port.JournalService, prev Prev) *CloseJournalModel {
	m := &CloseJournalModel{controller: controller, journalService: journalService, prev: prev}
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

func (m *CloseJournalModel) Init() tea.Cmd {
	return nil
}

func (m *CloseJournalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case closeJournalYes:
				if err := m.journalService.CloseJournal(); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.prev()
			case closeJournalNo:
				return m.prev()
			}
		}
	}
	return m, nil
}

func (m *CloseJournalModel) View() string {
	var s strings.Builder

	s.WriteString("  Закрыть журнал?\n\n")
	s.WriteString(m.menu.GenerateMenu())

	return s.String()
}
