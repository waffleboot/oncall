package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
	"github.com/waffleboot/oncall/pkg/tea/tabs"
)

func (m *TeaModel) updateVM(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
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
			if m.submitVM.Focused() {
				return m.exitScreen()
			}
		case "enter":
			if m.textinputVmName.Focused() || m.textinputVmNode.Focused() {
				var ok bool
				if m.tabs, cmd, ok = m.tabs.Next(); ok {
					return m, cmd
				}
			}
		}
	case button.PressedMsg:
		return m.runAndExitScreen(func() error {
			m.selectedVM.Name = m.textinputVmName.Value()
			m.selectedVM.Node = m.textinputVmNode.Value()
			m.selectedVM.Description = m.textinputVmDescription.Value()
			m.selectedItem.UpdateVM(m.selectedVM)
			if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
				return fmt.Errorf("update item: %w", err)
			}
			return nil
		})
	case string:
		if msg == "exit" {
			m.currentScreen = screenVMs
			return m, m.getItem
		}
	}

	switch {
	case m.textinputVmName.Focused():
		m.textinputVmName, cmd = m.textinputVmName.Update(msg)
		return m, cmd
	case m.textinputVmNode.Focused():
		m.textinputVmNode, cmd = m.textinputVmNode.Update(msg)
		return m, cmd
	case m.textinputVmDescription.Focused():
		m.textinputVmDescription, cmd = m.textinputVmDescription.Update(msg)
		return m, cmd
	case m.submitVM.Focused():
		m.submitVM, cmd = m.submitVM.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewVM() string {
	var s strings.Builder

	if m.selectedVM.Exists() {
		s.WriteString(fmt.Sprintf("ID: %d\n", m.selectedVM.ID))
	}
	s.WriteString("Name:\n  ")
	s.WriteString(m.textinputVmName.View())
	s.WriteString("\nNode:\n  ")
	s.WriteString(m.textinputVmNode.View())
	s.WriteString("\n")
	s.WriteString("Description:\n")
	s.WriteString(m.textinputVmDescription.View())
	s.WriteString("\n")
	s.WriteString(m.submitVM.View())
	s.WriteString("\n")

	return s.String()
}

func (m *TeaModel) resetVM() {
	m.textinputVmName = textinput.New()
	m.textinputVmName.Placeholder = "name"
	m.textinputVmName.Prompt = ""
	m.textinputVmName.Focus()
	m.textinputVmName.Width = 80
	m.textinputVmName.CharLimit = 1000
	m.textinputVmName.SetValue(m.selectedVM.Name)

	m.textinputVmNode = textinput.New()
	m.textinputVmNode.Placeholder = "node"
	m.textinputVmNode.Prompt = ""
	m.textinputVmNode.Blur()
	m.textinputVmNode.Width = 80
	m.textinputVmNode.CharLimit = 1000
	m.textinputVmNode.SetValue(m.selectedVM.Node)

	m.textinputVmDescription = textarea.New()
	m.textinputVmDescription.Placeholder = "vm description"
	m.textinputVmDescription.Blur()
	m.textinputVmDescription.Prompt = "  "
	m.textinputVmDescription.ShowLineNumbers = false
	m.textinputVmDescription.SetHeight(4)
	m.textinputVmDescription.SetWidth(80)
	m.textinputVmDescription.CharLimit = 1000
	m.textinputVmDescription.SetValue(m.selectedVM.Description)

	m.submitVM = button.New("submit")
	m.submitVM.Blur()

	m.tabs = tabs.New()
	m.tabs.Items = []tabs.Item{
		&m.textinputVmName,
		&m.textinputVmNode,
		&m.textinputVmDescription,
		&m.submitVM,
	}
	m.tabs.CanUp = func(tab int) bool {
		if tab == 2 && m.textinputVmDescription.Value() != "" {
			return false
		}
		return true
	}
	m.tabs.CanDown = m.tabs.CanUp
}
