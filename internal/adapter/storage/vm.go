package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type vm struct {
	ID          int       `json:"id"`
	Name        string    `json:"name,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (s *vm) fromDomain(vm model.VM) {
	s.ID = vm.ID
	s.Name = vm.Name
	s.DeletedAt = vm.DeletedAt
	s.Description = vm.Description
}

func (s *vm) toDomain() model.VM {
	return model.VM{
		ID:          s.ID,
		Name:        s.Name,
		DeletedAt:   s.DeletedAt,
		Description: s.Description,
	}
}
