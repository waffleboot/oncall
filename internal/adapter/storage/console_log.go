package storage

import "github.com/waffleboot/oncall/internal/model"

type consoleLog struct {
	ID       int
	VMID     string
	Filepath string
}

func (s *consoleLog) fromDomain(log model.ConsoleLog) {
	s.ID = log.ID
	s.VMID = log.VMID
	s.Filepath = log.Filepath
}

func (s *consoleLog) toDomain() model.ConsoleLog {
	return model.ConsoleLog{
		ID:       s.ID,
		VMID:     s.VMID,
		Filepath: s.Filepath,
	}
}
