package model

import (
	"slices"
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
	return s.AddedAt.Format(time.DateTime)
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
