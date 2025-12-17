package button

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	Model struct {
		value   string
		focused bool
	}
	PressedMsg struct {
		Value string
	}
)

func New(value string) Model { return Model{value: value} }

func (b Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !b.focused {
		return b, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return b, func() tea.Msg { return PressedMsg{Value: b.value} }
		}
	}

	return b, nil
}

func (b Model) View() string {
	if b.focused {
		return fmt.Sprintf("[[ %s ]]", strings.ToUpper(b.value))
	}
	return fmt.Sprintf("[ %s ]", strings.ToLower(b.value))
}

func (b Model) Focused() bool {
	return b.focused
}

func (b *Model) Focus() tea.Cmd {
	b.focused = true
	return nil
}

func (b *Model) Blur() {
	b.focused = false
}
