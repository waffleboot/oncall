package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type link struct {
	ID          int        `json:"id"`
	Link        string     `json:"link,omitempty"`
	Public      bool       `json:"public"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Description string     `json:"description,omitempty"`
}

func (s *link) fromDomain(link model.Link) {
	s.ID = link.ID
	s.Link = link.Address
	s.Public = link.Public
	s.DeletedAt = from(link.DeletedAt)
	s.Description = link.Description
}

func (s *link) toDomain() model.Link {
	return model.Link{
		ID:          s.ID,
		Address:     s.Link,
		Public:      s.Public,
		DeletedAt:   to(s.DeletedAt),
		Description: s.Description,
	}
}
