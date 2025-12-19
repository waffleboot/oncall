package model

import (
	"bufio"
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

func (n *Note) Exists() bool {
	return n.ID != 0
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

	fmt.Fprintf(sb, "#%d - %s", n.ID, n.trimNote())

	return sb.String()
}

func (n *Note) trimNote() string {
	s := bufio.NewScanner(strings.NewReader(n.Text))
	for s.Scan() {
		switch runes := []rune(strings.TrimSpace(s.Text())); {
		case len(runes) == 0:
			continue
		case len(runes) < 50:
			return string(runes)
		default:
			return string(runes[:50]) + " ..."
		}
	}
	if err := s.Err(); err != nil {
		return err.Error()
	}
	return "empty"
}

func (n *Note) ToPrint() string {
	return n.Text
}

func (s *Item) CreateNote() Note {
	return Note{Public: true}
}

func (s *Item) UpdateNote(note Note) {
	var maxID int
	var found bool
	for i, n := range s.Notes {
		if n.ID == note.ID {
			s.Notes[i] = note
			return
		}
		if n.Text == note.Text {
			found = true
		}
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	if found {
		return
	}
	note.ID = maxID + 1
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
