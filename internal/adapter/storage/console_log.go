package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type consoleLog struct {
	ID        int        `json:"id"`
	VMID      string     `json:"vmid"`
	FileID    string     `json:"file"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (s *consoleLog) fromDomain(log model.ConsoleLog) {
	s.ID = log.ID
	s.VMID = log.VMID
	s.FileID = log.FileID
	s.UpdatedAt = log.UpdatedAt.UTC()
	s.DeletedAt = from(log.DeletedAt)
}

func (s *consoleLog) toDomain() model.ConsoleLog {
	return model.ConsoleLog{
		ID:        s.ID,
		VMID:      s.VMID,
		FileID:    s.FileID,
		UpdatedAt: s.UpdatedAt,
		DeletedAt: to(s.DeletedAt),
	}
}
