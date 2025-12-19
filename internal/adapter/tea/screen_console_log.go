package tea

import (
	"fmt"
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
		case "up", "shift+tab":
			if m.textinputConsoleLogVMID.Focused() {
				m.textinputConsoleLogVMID.Blur()
				m.downloadConsoleLog.Focus()
				return m, nil
			}
			if m.textinputConsoleLogPath.Focused() {
				m.textinputConsoleLogPath.Blur()
				m.textinputConsoleLogVMID.Focus()
				return m, nil
			}
			if m.submitConsoleLog.Focused() {
				m.submitConsoleLog.Blur()
				m.textinputConsoleLogPath.Focus()
				return m, nil
			}
			if m.downloadConsoleLog.Focused() {
				m.downloadConsoleLog.Blur()
				m.submitConsoleLog.Focus()
				return m, nil
			}
		case "down", "tab":
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
				m.submitConsoleLog.Blur()
				m.downloadConsoleLog.Focus()
				return m, nil
			}
			if m.downloadConsoleLog.Focused() {
				m.downloadConsoleLog.Blur()
				m.textinputConsoleLogVMID.Focus()
				return m, nil
			}
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
		}
	case button.PressedMsg:
		m.selectedConsoleLog.VMID = m.textinputConsoleLogVMID.Value()
		m.selectedConsoleLog.Filepath = m.textinputConsoleLogPath.Value()
		return m.runAndExitScreen(func() error {
			m.selectedItem.UpdateConsoleLog(m.selectedConsoleLog)
			if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
				return fmt.Errorf("update item: %w", err)
			}
			return nil
		})
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
	sb.WriteString("\n\n")
	sb.WriteString(m.downloadConsoleLog.View())
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

	m.submitConsoleLog = button.New("submit")
	m.submitConsoleLog.Blur()

	filename := m.selectedConsoleLog.AddedAt.Format("2006-01-02-150405")
	if m.selectedConsoleLog.VMID != "" {
		filename = filename + "_" + m.selectedConsoleLog.VMID
	}

	m.downloadConsoleLog = button.New(fmt.Sprintf("download as %s.txt", filename))
	m.downloadConsoleLog.Blur()

	m.menuConsoleLogVMs.ResetMenu()
	m.menuConsoleLogVMs.AddGroupWithItems("vms", len(m.selectedItem.ActiveVMs()))
}
