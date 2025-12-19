package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type consoleLog struct {
	ID        int        `json:"id"`
	VMID      string     `json:"vmid"`
	Filepath  string     `json:"filepath"`
	AddedAt   time.Time  `json:"added_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (s *consoleLog) fromDomain(log model.ConsoleLog) {
	s.ID = log.ID
	s.VMID = log.VMID
	s.Filepath = log.FileID
	s.AddedAt = log.AddedAt.UTC()
	s.DeletedAt = from(log.DeletedAt)
}

func (s *consoleLog) toDomain() model.ConsoleLog {
	return model.ConsoleLog{
		ID:        s.ID,
		VMID:      s.VMID,
		FileID:    s.Filepath,
		AddedAt:   s.AddedAt,
		DeletedAt: to(s.DeletedAt),
	}
}
