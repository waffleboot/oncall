package model

import (
	"slices"
	"strings"
	"time"
)

type ConsoleLog struct {
	ID          int
	VMID        string
	FileID      string
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Description string
}

func (s ConsoleLog) HasFile() bool {
	return s.FileID != ""
}

func (s ConsoleLog) DownloadAs() string {
	filename := s.UpdatedAt.Format("2006-01-02-150405")
	if s.VMID != "" {
		filename = filename + "-" + s.VMID
	}
	return filename + ".txt"
}

func (s ConsoleLog) MenuItem() string {
	var sb strings.Builder
	sb.WriteString(s.UpdatedAt.Format(time.DateTime))
	if s.VMID != "" {
		sb.WriteString(" - ")
		sb.WriteString(s.VMID)
	}
	return sb.String()
}

func (s ConsoleLog) NotDeleted() bool {
	return s.DeletedAt.IsZero()
}

func (s *Item) CreateConsoleLog() ConsoleLog {
	return ConsoleLog{UpdatedAt: time.Now()}
}

func (s *Item) ActiveConsoleLogs() []ConsoleLog {
	consoleLogs := make([]ConsoleLog, 0, len(s.ConsoleLogs))
	for _, consoleLog := range s.ConsoleLogs {
		if consoleLog.NotDeleted() {
			consoleLogs = append(consoleLogs, consoleLog)
		}
	}
	return consoleLogs
}

func (s *Item) DeleteConsoleLog(consoleLog ConsoleLog) {
	if !consoleLog.HasFile() {
		s.ConsoleLogs = slices.DeleteFunc(s.ConsoleLogs, func(it ConsoleLog) bool {
			return it.ID == consoleLog.ID
		})
		return
	}
	for i := range s.ConsoleLogs {
		if s.ConsoleLogs[i].ID == consoleLog.ID {
			s.ConsoleLogs[i].DeletedAt = time.Now()
			return
		}
	}
}

func (s *Item) UpdateConsoleLog(consoleLog ConsoleLog) {
	consoleLog.UpdatedAt = time.Now()

	var maxID int
	for i, log := range s.ConsoleLogs {
		if log.ID == consoleLog.ID {
			s.ConsoleLogs[i] = consoleLog
			return
		}
		if log.ID > maxID {
			maxID = log.ID
		}
	}

	consoleLog.ID = maxID + 1
	s.ConsoleLogs = append(s.ConsoleLogs, consoleLog)
}
