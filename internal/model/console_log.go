package model

import (
	"slices"
	"strings"
	"time"
)

type ConsoleLog struct {
	ID          int
	VMID        string
	Filepath    string
	AddedAt     time.Time
	DeletedAt   time.Time
	Description string
}

func (s ConsoleLog) Empty() bool {
	return s.Filepath == ""
}

func (s ConsoleLog) MenuItem() string {
	var sb strings.Builder
	sb.WriteString(s.AddedAt.Format(time.DateTime))
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
	return ConsoleLog{AddedAt: time.Now()}
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

func (s *Item) DeleteConsoleLog(log ConsoleLog) {
	if log.Empty() {
		s.ConsoleLogs = slices.DeleteFunc(s.ConsoleLogs, func(it ConsoleLog) bool {
			return it.ID == log.ID
		})
		return
	}
	for i := range s.ConsoleLogs {
		if s.ConsoleLogs[i].ID == log.ID {
			s.ConsoleLogs[i].DeletedAt = time.Now()
			return
		}
	}
}

func (s *Item) UpdateConsoleLog(consoleLog ConsoleLog) {
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
