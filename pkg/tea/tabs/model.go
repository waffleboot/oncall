package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	Item interface {
		Blur()
		Focus() tea.Cmd
		Focused() bool
		View() string
	}
	Model struct {
		CanUp   func(tab int) bool
		CanDown func(tab int) bool
		Visible func(tab int) bool
		Items   []Item
		focus   int
	}
)

func New() Model {
	return Model{
		CanUp:   func(tab int) bool { return true },
		CanDown: func(tab int) bool { return true },
		Visible: func(tab int) bool { return true },
	}
}

func (b Model) Update(msg tea.Msg) (Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			b.Items[b.focus].Blur()
			b.focus = b.nextFocus(b.focus, +1)
			return b, b.Items[b.focus].Focus(), true
		case "shift+tab":
			b.Items[b.focus].Blur()
			b.focus = b.nextFocus(b.focus, -1)
			return b, b.Items[b.focus].Focus(), true
		case "up":
			if b.CanUp(b.focus) {
				b.Items[b.focus].Blur()
				b.focus = b.nextFocus(b.focus, -1)
				return b, b.Items[b.focus].Focus(), true
			}
		case "down":
			if b.CanDown(b.focus) {
				b.Items[b.focus].Blur()
				b.focus = b.nextFocus(b.focus, +1)
				return b, b.Items[b.focus].Focus(), true
			}
		}
	}
	return b, nil, false
}

func (b Model) Next() (Model, tea.Cmd, bool) {
	if b.CanDown(b.focus) {
		b.Items[b.focus].Blur()
		b.focus = b.nextFocus(b.focus, +1)
		return b, b.Items[b.focus].Focus(), true
	}
	return b, nil, false
}

func (b Model) nextFocus(current, delta int) int {
	current += delta
	if current >= 0 {
		current = current % len(b.Items)
	} else {
		current = len(b.Items) - 1
	}
	if b.Visible(current) {
		return current
	} else {
		return b.nextFocus(current, delta)
	}
}

//func (b Model) View() string {
//	var sb strings.Builder
//	for i := range b.Items {
//		if b.Hidden(i) {
//			continue
//		}
//		sb.WriteString(b.Items[i].View())
//		sb.WriteString("\n")
//	}
//	return sb.String()
//}
