package tea

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
)

func (m *TeaModel) updateConsoleLog(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "q":
			return m.exitScreen()
		case "enter":
			if m.textinputConsoleLogVMID.Focused() {
				m.textinputConsoleLogVMID.Blur()
				m.textinputConsoleLogPath.Focus()
				return m, nil
			}
			if m.textinputConsoleLogPath.Focused() {
				m.textinputConsoleLogPath.Blur()
				m.submitConsoleLog.Focus()
				return m, nil
			}
			if m.submitConsoleLog.Focused() {
				return m.exitScreen()
			}
		}
	case string:
		if msg == "exit" {
			m.currentScreen = screenConsoleLogs
			return m, m.getItem
		}
	}

	switch {
	case m.textinputConsoleLogVMID.Focused():
		m.textinputConsoleLogVMID, cmd = m.textinputConsoleLogVMID.Update(msg)
		return m, cmd
	case m.textinputConsoleLogPath.Focused():
		m.textinputConsoleLogPath, cmd = m.textinputConsoleLogPath.Update(msg)
		return m, cmd
	case m.submitConsoleLog.Focused():
		m.submitConsoleLog, cmd = m.submitConsoleLog.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewConsoleLog() string {
	var sb strings.Builder
	sb.WriteString("VMID:\n")
	sb.WriteString(m.textinputConsoleLogVMID.View())
	if len(m.selectedItem.ActiveVMs()) > 0 {
		sb.WriteString("\n\nVMID:\n")
		sb.WriteString(m.menuConsoleLogVMs.View())
	}
	sb.WriteString("\n\nFilepath:\n")
	sb.WriteString(m.textinputConsoleLogPath.View())
	sb.WriteString("\n\n")
	sb.WriteString(m.submitConsoleLog.View())
	return sb.String()
}

func (m *TeaModel) resetConsoleLog() {
	m.textinputConsoleLogVMID = textinput.New()
	m.textinputConsoleLogVMID.Placeholder = "vmid"
	m.textinputConsoleLogVMID.Prompt = ""
	m.textinputConsoleLogVMID.Focus()
	m.textinputConsoleLogVMID.Width = 80
	m.textinputConsoleLogVMID.CharLimit = 1000
	m.textinputConsoleLogVMID.SetValue(m.selectedConsoleLog.VMID)

	m.textinputConsoleLogPath = textinput.New()
	m.textinputConsoleLogPath.Placeholder = "filepath"
	m.textinputConsoleLogPath.Prompt = ""
	m.textinputConsoleLogPath.Blur()
	m.textinputConsoleLogPath.Width = 80
	m.textinputConsoleLogPath.CharLimit = 1000

	m.submitConsoleLog = button.New("Submit")
	m.submitConsoleLog.Blur()

	m.menuConsoleLogVMs.ResetMenu()
	m.menuConsoleLogVMs.AddGroupWithItems("vms", len(m.selectedItem.ActiveVMs()))
}
