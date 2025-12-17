package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type node struct {
	ID        int       `json:"id"`
	Name      string    `json:"name,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (s *node) fromDomain(node model.Node) {
	s.ID = node.ID
	s.Name = node.Name
	s.DeletedAt = node.DeletedAt.UTC()
}

func (s *node) toDomain() model.Node {
	return model.Node{
		ID:        s.ID,
		Name:      s.Name,
		DeletedAt: s.DeletedAt,
	}
}
