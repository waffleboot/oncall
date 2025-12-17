package model

import (
	"strings"
	"time"
)

type Note struct {
	ID        int
	Text      string
	DeletedAt time.Time
}

func (n *Note) IsDeleted() bool {
	return !n.DeletedAt.IsZero()
}

func (n *Note) MenuItem() string {
	switch {
	case n.Text == "":
		return "empty"
	case len(n.Text) < 50:
		return n.Text
	default:
		return n.Text[:50] + " ..."
	}
}

func (n *Note) ToPublish() string {
	return n.Text
}

func (s *Item) CreateNote() Note {
	var maxID int
	for i := range s.Notes {
		note := s.Notes[i]
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	note := Note{ID: maxID + 1}
	s.Notes = append(s.Notes, note)
	return note
}

func (s *Item) UpdateNote(note Note) {
	for i := range s.Notes {
		if s.Notes[i].ID == note.ID {
			s.Notes[i] = note
			break
		}
	}
}

func (s *Item) DeleteNote(note Note, at time.Time) {
	for i := range s.Notes {
		if s.Notes[i].ID == note.ID {
			s.Notes[i].DeletedAt = at
			break
		}
	}
}

func (s *Item) ActiveNotes() []Note {
	notes := make([]Note, 0, len(s.Notes))
	for _, note := range s.Notes {
		if !note.IsDeleted() {
			notes = append(notes, note)
		}
	}
	return notes
}

func (n *Note) Printed() bool {
	return !n.IsDeleted() && strings.TrimSpace(n.Text) != ""
}

func (s *Item) PrintedNotes() []Note {
	notes := make([]Note, 0, len(s.Notes))
	for _, note := range s.Notes {
		if note.Printed() {
			notes = append(notes, note)
		}
	}
	return notes
}
