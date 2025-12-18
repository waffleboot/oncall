package menu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

const menuDelimiter = "---"

type (
	Model struct {
		labelGen    func(group string, pos int) string
		groupNames  []string
		groupSizes  []int
		cursor      int
		selected    string
		notSelected string
		log         *zap.Logger
	}
	MenuOption func(*Model)
)

func WithLogger(log *zap.Logger) MenuOption {
	return func(menu *Model) {
		menu.log = log
	}
}

func WithSelection(selected, notSelected string) MenuOption {
	return func(menu *Model) {
		menu.selected = selected
		menu.notSelected = notSelected
	}
}

func New(labelGen func(group string, pos int) string, opts ...MenuOption) *Model {
	m := &Model{labelGen: labelGen}
	m.log = zap.NewNop()
	m.selected = ">"
	m.notSelected = " "
	for i := range opts {
		opts[i](m)
	}
	return m
}

func (m *Model) ResetCursor() {
	m.cursor = 0
}

func (m *Model) ResetMenu() {
	m.groupNames = nil
	m.groupSizes = nil
}

func (m *Model) AddGroup(group string) {
	m.groupNames = append(m.groupNames, group)
	m.groupSizes = append(m.groupSizes, 1)
}

func (m *Model) AddGroupWithItems(group string, itemsCount int) {
	m.groupNames = append(m.groupNames, group)
	m.groupSizes = append(m.groupSizes, itemsCount)
}

func (m *Model) AdjustCursor() {
	if g, _ := m.GetGroup(); g == "" {
		m.cursor = m.maxCursor()
	}
}

func (m *Model) AddDelimiter() {
	m.groupNames = append(m.groupNames, menuDelimiter)
	m.groupSizes = append(m.groupSizes, 1)
}

func (m *Model) View() string {
	var cursor int
	var s strings.Builder

	for i, group := range m.groupNames {
		if group == menuDelimiter {
			s.WriteString("\n")
			continue
		}
		for pos := 0; pos < m.groupSizes[i]; pos++ {
			if cursor == m.cursor {
				s.WriteString(m.selected)

			} else {
				s.WriteString(m.notSelected)
			}
			s.WriteString(" ")

			label := m.labelGen(group, pos)
			if label == "" {
				label = group
			}

			s.WriteString(label)
			s.WriteString("\n")
			cursor++
		}
	}

	m.log.Info("view", zap.String("view", s.String()))

	return s.String()
}

func (m *Model) GetGroup() (group string, pos int) {
	var count int

	for i, group := range m.groupNames {
		if group == menuDelimiter {
			continue
		}
		a := count
		b := count + m.groupSizes[i]
		if a <= m.cursor && m.cursor < b {
			return group, m.cursor - a
		}
		count = b
	}

	return "", 0
}

func (m *Model) maxCursor() int {
	var count int
	for i, group := range m.groupNames {
		if group == menuDelimiter {
			continue
		}
		count += m.groupSizes[i]
	}
	return count - 1
}

func (m *Model) MoveCursorUp() {
	if m.cursor > 0 {
		m.cursor = m.cursor - 1
	} else {
		m.cursor = m.maxCursor()
	}
}

func (m *Model) MoveCursorDown() {
	if m.cursor < m.maxCursor() {
		m.cursor = m.cursor + 1
	} else {
		m.cursor = 0
	}
}

func (m *Model) JumpToItem(group string, f func(pos int) (found bool)) {
	m.jumpTo(group, func(pos int) bool {
		return f(pos)
	})
}

func (m *Model) JumpToPos(toGroup string, pos int) {
	var cursor int
	for i, group := range m.groupNames {
		if group == menuDelimiter {
			continue
		}
		if group != toGroup {
			cursor += m.groupSizes[i]
			continue
		}
		if pos < m.groupSizes[i] {
			m.cursor = cursor + pos
		} else {
			return
		}
	}
}

func (m *Model) JumpToGroup(group string) {
	m.jumpTo(group, func(_ int) bool {
		return true
	})
}

func (m *Model) jumpTo(toGroup string, find func(pos int) bool) {
	var cursor int

	log := m.log.With(zap.String("toGroup", toGroup), zap.Int("groups", len(m.groupNames)))

	log.Debug("jump to")

	for i, group := range m.groupNames {
		if group == menuDelimiter {
			log.Debug(
				"skip delimiter",
				zap.Int("pos", i),
				zap.Int("cursor", cursor),
			)
			continue
		}
		if group != toGroup {
			cursor += m.groupSizes[i]

			log.Debug(
				"skip group",
				zap.String("group", group),
				zap.Int("count", m.groupSizes[i]),
				zap.Int("cursor", cursor),
			)

			continue
		}

		log.Debug(
			"group found",
			zap.String("group", group),
			zap.Int("size", m.groupSizes[i]),
			zap.Int("cursor", cursor),
		)

		for pos := 0; pos < m.groupSizes[i]; pos++ {
			found := find(pos)

			log.Debug(
				"check pos",
				zap.String("group", group),
				zap.Int("size", m.groupSizes[i]),
				zap.Int("cursor", cursor),
				zap.Bool("found", found),
				zap.Int("pos", pos),
			)

			if found {
				m.cursor = cursor
				return
			}
			cursor++
		}
	}
}

func (m *Model) ProcessMsg(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.MoveCursorUp()
			return true
		case "down", "j":
			m.MoveCursorDown()
			return true
		}
	}
	return false
}
