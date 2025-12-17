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

func (t Model) Update(msg tea.Msg) (Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			t.Items[t.focus].Blur()
			t.focus = t.nextFocus(t.focus, +1)
			return t, t.Items[t.focus].Focus(), true
		case "shift+tab":
			t.Items[t.focus].Blur()
			t.focus = t.nextFocus(t.focus, -1)
			return t, t.Items[t.focus].Focus(), true
		case "up":
			if t.CanUp(t.focus) {
				t.Items[t.focus].Blur()
				t.focus = t.nextFocus(t.focus, -1)
				return t, t.Items[t.focus].Focus(), true
			}
		case "down":
			if t.CanDown(t.focus) {
				t.Items[t.focus].Blur()
				t.focus = t.nextFocus(t.focus, +1)
				return t, t.Items[t.focus].Focus(), true
			}
		}
	}
	return t, nil, false
}

func (t Model) Next() (Model, tea.Cmd, bool) {
	if t.CanDown(t.focus) {
		t.Items[t.focus].Blur()
		t.focus = t.nextFocus(t.focus, +1)
		return t, t.Items[t.focus].Focus(), true
	}
	return t, nil, false
}

func (t Model) nextFocus(current, delta int) int {
	current += delta
	if current >= 0 {
		current = current % len(t.Items)
	} else {
		current = len(t.Items) - 1
	}
	if t.Visible(current) {
		return current
	} else {
		return t.nextFocus(current, delta)
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
