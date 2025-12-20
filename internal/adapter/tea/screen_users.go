package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateUsers(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuUsers.Update(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "e":
			return m.exitScreen()
		case "enter", " ":
			switch g, p := m.menuUsers.GetGroup(); g {
			case "exit":
				return m.exitScreen()
			case "users":
				return m, func() tea.Msg {
					if err := m.userService.SetUser(m.users[p]); err != nil {
						return fmt.Errorf("set user: %w", err)
					}
					return "exit"
				}
			}
		}
	case string:
		if msg == "exit" {
			m.currentScreen = screenItems
			return m, nil
		}
	}

	return m, nil
}

func (m *TeaModel) viewUsers() string {
	return m.menuUsers.View()
}

func (m *TeaModel) resetUsers(user *model.User) {
	m.menuUsers.ResetMenu()
	m.menuUsers.AddGroup("exit")
	m.menuUsers.AddGroupWithItems("users", len(m.users))
	if user != nil {
		m.menuUsers.JumpToItem("users", func(pos int) (found bool) {
			return m.users[pos].Nick == user.Nick
		})
	} else {
		m.menuUsers.JumpToGroup("exit")
	}
}
