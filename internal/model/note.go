package model

import (
	"fmt"
	"strings"
	"time"
)

type Note struct {
	ID        int
	Text      string
	Public    bool
	DeletedAt time.Time
}

func (n *Note) IsDeleted() bool {
	return !n.DeletedAt.IsZero()
}

func (n *Note) MenuItem() string {
	sb := new(strings.Builder)

	if n.Public {
		sb.WriteString("  ")
	} else {
		sb.WriteString("p ")
	}

	fmt.Fprintf(sb, "#%d - ", n.ID)

	switch {
	case n.Text == "":
		sb.WriteString("empty")
	case len(n.Text) < 50:
		sb.WriteString(n.Text)
	default:
		sb.WriteString(n.Text[:50] + " ...")
	}

	return sb.String()
}

func (n *Note) ToPrint() string {
	return n.Text
}

func (s *Item) CreateNote() Note {
	var maxID int
	for _, note := range s.Notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return Note{ID: maxID + 1, Public: true}
}

func (s *Item) UpdateNote(note Note) {
	for i := range s.Notes {
		if s.Notes[i].ID == note.ID {
			s.Notes[i] = note
			return
		}
	}
	s.Notes = append(s.Notes, note)
}

func (s *Item) DeleteNote(note Note) {
	for i := range s.Notes {
		if s.Notes[i].ID == note.ID {
			s.Notes[i].DeletedAt = time.Now()
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
	return !n.IsDeleted() && n.Public && strings.TrimSpace(n.Text) != ""
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
