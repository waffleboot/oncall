package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type vm struct {
	ID          int       `json:"id"`
	Name        string    `json:"name,omitempty"`
	Node        string    `json:"node,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (s *vm) fromDomain(vm model.VM) {
	s.ID = vm.ID
	s.Name = vm.Name
	s.Node = vm.Node
	s.DeletedAt = vm.DeletedAt.UTC()
	s.Description = vm.Description
}

func (s *vm) toDomain() model.VM {
	return model.VM{
		ID:          s.ID,
		Name:        s.Name,
		Node:        s.Node,
		DeletedAt:   s.DeletedAt,
		Description: s.Description,
	}
}
