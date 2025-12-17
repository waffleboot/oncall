package tea

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
	"github.com/waffleboot/oncall/pkg/tea/tabs"
)

func (m *TeaModel) updateNewNodes(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	var ok bool
	if m.tabs, cmd, ok = m.tabs.Update(msg); ok {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "q":
			if m.submitNodes.Focused() {
				return m.exitScreen()
			}
		}
	case button.PressedMsg:
		return m.runAndExitScreen(func() error {
			s := bufio.NewScanner(strings.NewReader(m.textinputNodes.Value()))
			for s.Scan() {
				t := strings.TrimSpace(s.Text())
				if t != "" {
					node := m.selectedItem.CreateNode()
					node.Name = t
					m.selectedItem.UpdateNode(node)
				}
			}
			if err := s.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
				return fmt.Errorf("update item: %w", err)
			}
			return nil
		})
	case string:
		switch msg {
		case "exit":
			m.currentScreen = screenNodes
			return m, m.getItem
		}
	}

	switch {
	case m.textinputNodes.Focused():
		m.textinputNodes, cmd = m.textinputNodes.Update(msg)
		return m, cmd
	case m.submitNodes.Focused():
		m.submitNodes, cmd = m.submitNodes.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m *TeaModel) viewNewNodes() string {
	var s strings.Builder
	s.WriteString("Nodes:\n")
	s.WriteString(m.textinputNodes.View())
	s.WriteString("\n")
	s.WriteString(m.submitNodes.View())
	s.WriteString("\n")
	return s.String()
}

func (m *TeaModel) resetNewNodes() {
	m.textinputNodes = textarea.New()
	m.textinputNodes.Placeholder = "nodes"
	m.textinputNodes.Focus()
	m.textinputNodes.ShowLineNumbers = false
	m.textinputNodes.SetHeight(4)
	m.textinputNodes.SetWidth(80)
	m.textinputNodes.CharLimit = 1000

	m.submitNodes = button.New("submit")
	m.submitNodes.Blur()

	m.tabs = tabs.New()
	m.tabs.Items = []tabs.Item{
		&m.textinputNodes,
		&m.submitNodes,
	}
	m.tabs.CanUp = func(tab int) bool {
		if tab == 0 && m.textinputNodes.Value() != "" {
			return false
		}
		return true
	}
	m.tabs.CanDown = m.tabs.CanUp
}
