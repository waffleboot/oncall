package storage

import (
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type item struct {
	ID          int          `json:"id"`
	SleepAt     time.Time    `json:"sleep_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
	DeletedAt   time.Time    `json:"deleted_at,omitempty"`
	ClosedAt    time.Time    `json:"closed_at,omitempty"`
	Links       []link       `json:"links,omitempty"`
	Notes       []note       `json:"notes,omitempty"`
	Nodes       []node       `json:"nodes,omitempty"`
	VMs         []vm         `json:"vms,omitempty"`
	Type        string       `json:"type,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	ConsoleLogs []consoleLog `json:"console_logs,omitempty"`
}

func (s *item) NotDeleted() bool {
	return s.DeletedAt.IsZero()
}

func (s *item) fromDomain(item model.Item) {
	s.ID = item.ID
	s.SleepAt = item.SleepAt.UTC()
	s.CreatedAt = item.CreatedAt.UTC()
	s.UpdatedAt = item.UpdatedAt.UTC()
	s.DeletedAt = item.DeletedAt.UTC()
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

	s.Links = make([]link, len(item.Links))
	for i := range item.Links {
		s.Links[i].fromDomain(item.Links[i])
	}

	s.Links = make([]link, len(item.Links))
	for i := range item.Links {
		s.Links[i].fromDomain(item.Links[i])
	}

	s.ConsoleLogs = make([]consoleLog, len(item.ConsoleLogs))
	for i := range item.ConsoleLogs {
		s.ConsoleLogs[i].fromDomain(item.ConsoleLogs[i])
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

	consoleLogs := make([]model.ConsoleLog, len(s.ConsoleLogs))
	for i := range s.ConsoleLogs {
		consoleLogs[i] = s.ConsoleLogs[i].toDomain()
	}

	return model.Item{
		ID:          s.ID,
		SleepAt:     s.SleepAt,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt,
		ClosedAt:    s.ClosedAt,
		Type:        model.ItemType(s.Type),
		Title:       s.Title,
		Description: s.Description,
		Links:       links,
		Nodes:       nodes,
		Notes:       notes,
		VMs:         vms,
		ConsoleLogs: consoleLogs,
	}
}
