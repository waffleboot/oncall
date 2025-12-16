package model

import (
	"time"
)

type VM struct {
	ID          int
	Name        string
	Node        string
	DeletedAt   time.Time
	Description string
}

func (s *VM) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}

func (s *Item) ActiveVMs() []VM {
	vms := make([]VM, 0, len(s.VMs))
	for _, vm := range s.VMs {
		if !vm.IsDeleted() {
			vms = append(vms, vm)
		}
	}
	return vms
}

func (s *Item) CreateVM() VM {
	var maxID int
	for i := range s.VMs {
		vm := s.VMs[i]
		if vm.ID > maxID {
			maxID = vm.ID
		}
	}
	vm := VM{ID: maxID + 1}
	s.VMs = append(s.VMs, vm)
	return vm
}

func (s *Item) UpdateVM(vm VM) {
	for i := range s.VMs {
		if s.VMs[i].ID == vm.ID {
			s.VMs[i] = vm
			break
		}
	}
}

func (s *Item) DeleteVM(vm VM, at time.Time) {
	for i := range s.VMs {
		if s.VMs[i].ID == vm.ID {
			s.VMs[i].DeletedAt = at
			break
		}
	}
}
