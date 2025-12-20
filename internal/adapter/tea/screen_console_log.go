package tea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/pkg/tea/button"
)

func (m *TeaModel) updateConsoleLog(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	m.consoleLogError = nil

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
				m.selectedConsoleLog.VMID = m.textinputConsoleLogVMID.Value()
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
		if msg.Value == " submit " {
			return m.runAndExitScreen(func() error {
				src := m.textinputConsoleLogPath.Value()
				if src != "" {
					fileID, err := m.fileStorage.UploadFile(src)
					if err != nil {
						return consoleLogErrorMsg(fmt.Errorf("upload file: %w", err))
					}
					m.selectedConsoleLog.FileID = fileID
				}
				m.selectedConsoleLog.VMID = m.textinputConsoleLogVMID.Value()

				m.selectedItem.UpdateConsoleLog(m.selectedConsoleLog)
				if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				}
				return nil
			})
		}
		return m.runAndExitScreen(func() error {
			fileID := m.selectedConsoleLog.FileID
			destination := m.downloadConsoleLogAs()
			if err := m.fileStorage.DownloadFile(fileID, destination); err != nil {
				return consoleLogErrorMsg(fmt.Errorf("download file: %w", err))
			}
			return nil
		})
	case consoleLogErrorMsg:
		m.consoleLogError = msg
		return m, nil
	case string:
		if msg == "exit" {
			m.currentScreen = screenConsoleLogs
			return m, m.getItem
		}
	}

	switch {
	case m.textinputConsoleLogVMID.Focused():
		m.textinputConsoleLogVMID, cmd = m.textinputConsoleLogVMID.Update(msg)
		m.selectedConsoleLog.VMID = m.textinputConsoleLogVMID.Value()
		if m.selectedConsoleLog.HasFile() {
			m.textinputConsoleLogPath.Placeholder = m.downloadConsoleLogAs()
		}
		return m, cmd
	case m.textinputConsoleLogPath.Focused():
		m.textinputConsoleLogPath, cmd = m.textinputConsoleLogPath.Update(msg)
		return m, cmd
	case m.submitConsoleLog.Focused():
		m.submitConsoleLog, cmd = m.submitConsoleLog.Update(msg)
		return m, cmd
	case m.downloadConsoleLog.Focused():
		m.downloadConsoleLog, cmd = m.downloadConsoleLog.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *TeaModel) viewConsoleLog() string {
	var sb strings.Builder
	sb.WriteString("VMID: ")
	sb.WriteString(m.textinputConsoleLogVMID.View())
	// if len(m.selectedItem.ActiveVMs()) > 0 {
	// 	sb.WriteString("\n\nVMID:\n")
	// 	sb.WriteString(m.menuConsoleLogVMs.View())
	// }
	sb.WriteString("\n\nFilepath to upload or download: ")
	sb.WriteString(m.textinputConsoleLogPath.View())
	sb.WriteString("\n\n")
	sb.WriteString(m.submitConsoleLog.View())
	if m.selectedConsoleLog.HasFile() {
		sb.WriteString("\n\n")
		sb.WriteString(m.downloadConsoleLog.View())
		sb.WriteString(" as ")
		sb.WriteString(m.downloadConsoleLogAs())
	}
	if m.consoleLogError != nil {
		sb.WriteString("\n\n")
		sb.WriteString(m.consoleLogError.Error())
	}
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
	m.textinputConsoleLogPath.Prompt = ""
	m.textinputConsoleLogPath.Blur()
	m.textinputConsoleLogPath.Width = 80
	m.textinputConsoleLogPath.CharLimit = 1000
	if m.selectedConsoleLog.HasFile() {
		m.textinputConsoleLogPath.Placeholder = m.downloadConsoleLogAs()
	} else {
		m.textinputConsoleLogPath.Placeholder = "filepath"
	}

	m.submitConsoleLog = button.New(" submit ")
	m.submitConsoleLog.Blur()

	m.downloadConsoleLog = button.New("download")
	m.downloadConsoleLog.Blur()

	m.menuConsoleLogVMs.ResetMenu()
	m.menuConsoleLogVMs.AddGroupWithItems("vms", len(m.selectedItem.ActiveVMs()))

	if len(m.selectedItem.ActiveVMs()) == 1 {
		m.textinputConsoleLogVMID.SetValue(m.selectedItem.ActiveVMs()[0].Name)
	}
}

func (m *TeaModel) downloadConsoleLogAs() string {
	if m.textinputConsoleLogPath.Value() != "" {
		return m.textinputConsoleLogPath.Value()
	} else {
		return m.selectedConsoleLog.DownloadAs()
	}
}
