package tea

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateConsoleLog(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "q":
			return m.exitScreen()
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
	}

	return m, nil
}

func (m *TeaModel) viewConsoleLog() string {
	var sb strings.Builder
	sb.WriteString("VMID:\n")
	sb.WriteString(m.textinputConsoleLogVMID.View())
	sb.WriteString("\nVMID:\n")
	sb.WriteString(m.menuConsoleLogVMs.View())
	sb.WriteString("\n")
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

	m.menuConsoleLogVMs.ResetMenu()
	m.menuConsoleLogVMs.AddGroupWithItems("vms", len(m.selectedItem.ActiveVMs()))
}
