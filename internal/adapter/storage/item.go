package storage

import "github.com/waffleboot/oncall/internal/model"

func (s *storedItem) fromDomain(item model.Item) {
	s.ID = item.ID
	s.Num = item.Num
	s.SleepAt = item.SleepAt.UTC()
	s.ClosedAt = item.ClosedAt.UTC()
	s.Type = string(item.Type)
	s.Title = item.Title
	s.Description = item.Description
	s.Links = make([]storedLink, len(item.Links))
	for i := range item.Links {
		s.Links[i].fromDomain(item.Links[i])
	}
	s.VMs = make([]vm, len(item.VMs))
	for i := range item.VMs {
		s.VMs[i].fromDomain(item.VMs[i])
	}
}

func (s *storedItem) toDomain() model.Item {
	vms := make([]model.VM, len(s.VMs))
	links := make([]model.Link, len(s.Links))
	for i := range s.Links {
		links[i] = s.Links[i].toDomain()
	}
	for i := range s.VMs {
		vms[i] = s.VMs[i].toDomain()
	}
	return model.Item{
		ID:          s.ID,
		Num:         s.Num,
		SleepAt:     s.SleepAt,
		ClosedAt:    s.ClosedAt,
		Type:        model.ItemType(s.Type),
		Title:       s.Title,
		Description: s.Description,
		Links:       links,
		VMs:         vms,
	}
}
