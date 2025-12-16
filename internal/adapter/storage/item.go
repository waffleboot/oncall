package storage

import "github.com/waffleboot/oncall/internal/model"

func (s *item) fromDomain(item model.Item) {
	s.ID = item.ID
	s.Num = item.Num
	s.SleepAt = item.SleepAt.UTC()
	s.ClosedAt = item.ClosedAt.UTC()
	s.Type = string(item.Type)
	s.Title = item.Title
	s.Description = item.Description

	s.VMs = make([]vm, len(item.VMs))
	for i := range item.VMs {
		s.VMs[i].fromDomain(item.VMs[i])
	}

	s.Notes = make([]note, len(item.Notes))
	for i := range item.Notes {
		s.Notes[i].fromDomain(item.Notes[i])
	}

	s.Nodes = make([]node, len(item.Nodes))
	for i := range item.Nodes {
		s.Nodes[i].fromDomain(item.Nodes[i])
	}

	s.Links = make([]storedLink, len(item.Links))
	for i := range item.Links {
		s.Links[i].fromDomain(item.Links[i])
	}
}

func (s *item) toDomain() model.Item {
	vms := make([]model.VM, len(s.VMs))
	for i := range s.VMs {
		vms[i] = s.VMs[i].toDomain()
	}

	nodes := make([]model.Node, len(s.Nodes))
	for i := range s.Nodes {
		nodes[i] = s.Nodes[i].toDomain()
	}

	notes := make([]model.Note, len(s.Notes))
	for i := range s.Notes {
		notes[i] = s.Notes[i].toDomain()
	}

	links := make([]model.Link, len(s.Links))
	for i := range s.Links {
		links[i] = s.Links[i].toDomain()
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
		Nodes:       nodes,
		Notes:       notes,
		VMs:         vms,
	}
}
