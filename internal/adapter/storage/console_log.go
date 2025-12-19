package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type consoleLog struct {
	ID        int
	VMID      string
	Filepath  string
	AddedAt   time.Time
	DeletedAt *time.Time
}

func (s *consoleLog) fromDomain(log model.ConsoleLog) {
	s.ID = log.ID
	s.VMID = log.VMID
	s.Filepath = log.Filepath
	s.AddedAt = log.AddedAt.UTC()
	s.DeletedAt = from(log.DeletedAt)
}

func (s *consoleLog) toDomain() model.ConsoleLog {
	return model.ConsoleLog{
		ID:        s.ID,
		VMID:      s.VMID,
		Filepath:  s.Filepath,
		AddedAt:   s.AddedAt,
		DeletedAt: to(s.DeletedAt),
	}
}
