package model

import (
	"fmt"
	"strings"
	"time"
)

type VM struct {
	ID          int
	Name        string
	Node        string
	DeletedAt   time.Time
	Description string
}

func (v *VM) Exists() bool {
	return v.ID != 0
}

func (v *VM) NotDeleted() bool {
	return v.DeletedAt.IsZero()
}

func (v *VM) HasNode() bool {
	return strings.TrimSpace(v.Node) != ""
}

func (v *VM) MenuItem() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("#%d - ", v.ID))

	if v.Name == "" {
		sb.WriteString("empty")
	} else {
		sb.WriteString(v.Name)
	}

	switch {
	case len(v.Description) > 50:
		sb.WriteString(" - ")
		sb.WriteString(v.Description[:50])
	case v.Description != "":
		sb.WriteString(" - ")
		sb.WriteString(v.Description)
	}

	return sb.String()
}

func (v *VM) ToPrint() string {
	if v.Node != "" {
		return fmt.Sprintf("vm: %s\nhost: %s", v.Name, v.Node)
	}
	return fmt.Sprintf("vm: %s", v.Name)
}

func (s *Item) ActiveVMs() []VM {
	vms := make([]VM, 0, len(s.VMs))
	for _, vm := range s.VMs {
		if vm.NotDeleted() {
			vms = append(vms, vm)
		}
	}
	return vms
}

func (s *Item) CreateVM() VM {
	return VM{}
}

func (s *Item) UpdateVM(vm VM) {
	var maxID int
	for i, v := range s.VMs {
		if v.ID == vm.ID {
			s.VMs[i] = vm
			return
		}
		if v.ID > maxID {
			maxID = v.ID
		}
	}
	vm.ID = maxID + 1
	s.VMs = append(s.VMs, vm)
}

func (s *Item) DeleteVM(vm VM) {
	for i := range s.VMs {
		if s.VMs[i].ID == vm.ID {
			s.VMs[i].DeletedAt = time.Now()
			break
		}
	}
}

func (v *VM) Printed() bool {
	return v.NotDeleted() && strings.TrimSpace(v.Name) != ""
}

func (s *Item) PrintedVMs() []VM {
	vms := make([]VM, 0, len(s.VMs))
	for _, vm := range s.VMs {
		if vm.Printed() {
			vms = append(vms, vm)
		}
	}
	return vms
}
