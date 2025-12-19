package model

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Link struct {
	ID          int
	Public      bool
	Address     string
	DeletedAt   time.Time
	Description string
}

func (s *Item) CreateLink() Link {
	return Link{Public: true}
}

func (s *Item) UpdateLink(link Link) {
	var maxID int
	var found bool
	for i, l := range s.Links {
		if l.ID == link.ID {
			s.Links[i] = link
			return
		}
		if l.Address == link.Address {
			found = true
		}
		if l.ID > maxID {
			maxID = l.ID
		}
	}
	if found {
		return
	}
	link.ID = maxID + 1
	s.Links = append(s.Links, link)
}

func (s *Item) DeleteLink(link Link) {
	if link.Empty() {
		s.Links = slices.DeleteFunc(s.Links, func(it Link) bool {
			return it.ID == link.ID
		})
		return
	}
	for i := range s.Links {
		if s.Links[i].ID == link.ID {
			s.Links[i].DeletedAt = time.Now()
			return
		}
	}
}

func (s *Link) HasID() bool {
	return s.ID != 0
}

func (s *Link) Empty() bool {
	return strings.TrimSpace(s.Address) == "" && strings.TrimSpace(s.Description) == ""
}

func (s *Link) ToPrint() string {
	description := strings.TrimSpace(s.Description)
	if description != "" {
		return fmt.Sprintf("%s - %s", s.Address, description)
	}
	return s.Address
}

func (s *Link) MenuItem() string {
	sb := new(strings.Builder)

	if s.Public {
		sb.WriteString("  ")
	} else {
		sb.WriteString("p ")
	}

	// fmt.Fprintf(sb, "#%d ", s.ID)

	if s.Address == "" {
		sb.WriteString("empty")
	} else {
		sb.WriteString(s.Address)
	}

	switch {
	case len(s.Description) > 50:
		sb.WriteString(" - ")
		sb.WriteString(s.Description[:50])
	case s.Description != "":
		sb.WriteString(" - ")
		sb.WriteString(s.Description)
	}

	return sb.String()
}

func (s *Link) NotDeleted() bool {
	return s.DeletedAt.IsZero()
}

func (s *Item) ActiveLinks() []Link {
	links := make([]Link, 0, len(s.Links))
	for _, link := range s.Links {
		if link.NotDeleted() {
			links = append(links, link)
		}
	}
	return links
}

func (s *Item) PrintedLinks() []Link {
	links := make([]Link, 0, len(s.Links))
	for _, link := range s.Links {
		if link.NotDeleted() && link.Public && strings.TrimSpace(link.Address) != "" {
			links = append(links, link)
		}
	}
	return links
}
