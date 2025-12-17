package tea

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateNewNodes(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.textInput {
			case "nodes":
				m.textInput = "submit"
				m.textinputNodes.Blur()
			case "submit":
				m.textInput = "nodes"
				m.textinputNodes.Focus()
			}
			return m, nil
		case "up":
			switch m.textInput {
			case "nodes":
				if len(m.textinputNodes.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNodes.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "nodes"
				m.textinputNodes.Focus()
				return m, nil
			}
		case "down":
			switch m.textInput {
			case "nodes":
				if len(m.textinputNodes.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNodes.Blur()
					return m, nil
				}
			case "submit":
				m.textInput = "nodes"
				m.textinputNodes.Focus()
				return m, nil
			}
		case "esc":
			m.currentScreen = screenNodes
			return m, m.getItem
		case "q":
			if m.textInput != "nodes" {
				m.currentScreen = screenNodes
				return m, m.getItem
			}
		case "enter":
			switch m.textInput {
			case "nodes":
				if len(m.textinputNodes.Value()) == 0 {
					m.textInput = "submit"
					m.textinputNodes.Blur()
					return m, nil
				}
			case "submit":
				return m, func() tea.Msg {
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
					return m.getItem()
				}
			}
		}
	case model.Item:
		m.currentScreen = screenNodes
		return m, m.getItem
	}

	m.textinputNodes, cmd = m.textinputNodes.Update(msg)

	return m, cmd
}

func (m *TeaModel) viewNewNodes() string {
	var s strings.Builder
	s.WriteString("Nodes:\n")
	s.WriteString(m.textinputNodes.View())
	s.WriteString("\n")
	if m.textInput == "submit" {
		s.WriteString("[[ SUBMIT ]]\n")
	} else {
		s.WriteString("[ submit ]\n")
	}
	return s.String()
}

func (m *TeaModel) resetNewNodes() {
	m.textInput = "nodes"
	m.textinputNodes = textarea.New()
	m.textinputNodes.Placeholder = "nodes"
	m.textinputNodes.Focus()
	m.textinputNodes.ShowLineNumbers = false
	m.textinputNodes.SetHeight(4)
	m.textinputNodes.SetWidth(80)
	m.textinputNodes.CharLimit = 1000
}
