package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type note struct {
	ID        int       `json:"id"`
	Text      string    `json:"text,omitempty"`
	Public    bool      `json:"public,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (s *note) fromDomain(note model.Note) {
	s.ID = note.ID
	s.Text = note.Text
	s.Public = note.Public
	s.DeletedAt = note.DeletedAt
}

func (s *note) toDomain() model.Note {
	return model.Note{
		ID:        s.ID,
		Text:      s.Text,
		Public:    s.Public,
		DeletedAt: s.DeletedAt,
	}
}
