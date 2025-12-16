package model

import (
	"cmp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ItemTypeInc   ItemType = "inc"
	ItemTypeAdhoc ItemType = "adhoc"
	ItemTypeAsk   ItemType = "ask"
	ItemTypeAlert ItemType = "alert"
)

type (
	ItemType string
	Item     struct {
		ID          uuid.UUID
		Num         int
		SleepAt     time.Time
		ClosedAt    time.Time
		Type        ItemType
		Links       []Link
		VMs         []VM
		Title       string
		Description string
	}
	Link struct {
		ID          int
		Public      bool
		Address     string
		DeletedAt   time.Time
		Description string
	}
)

func (i *Item) IsActive() bool {
	return !i.IsClosed() && i.SleepAt.IsZero()
}

func (i *Item) IsSleep() bool {
	return !i.IsClosed() && !i.SleepAt.IsZero()
}

func (i *Item) IsClosed() bool {
	return !i.ClosedAt.IsZero()
}

func (i *Item) Sleep(at time.Time) {
	i.SleepAt = at
}

func (i *Item) Awake() {
	i.SleepAt = time.Time{}
}

func (i *Item) Close(at time.Time) {
	i.ClosedAt = at
}

func (t ItemType) String() string {
	switch t {
	case ItemTypeInc:
		return "Инцидент"
	case ItemTypeAdhoc:
		return "Adhoc"
	case ItemTypeAsk:
		return "Обращение"
	case ItemTypeAlert:
		return "Alert"
	}
	return "Unknown"
}

func (t ItemType) Compare(o ItemType) int {
	pri := func(t ItemType) int {
		switch t {
		case ItemTypeAsk:
			return 1
		case ItemTypeInc:
			return 2
		case ItemTypeAlert:
			return 3
		case ItemTypeAdhoc:
			return 4
		default:
			return 5
		}
	}
	return cmp.Compare(pri(t), pri(o))
}

func (t Item) Compare(o Item) int {
	pri := func(t Item) int {
		switch {
		case t.IsActive():
			return 1
		case t.IsSleep():
			return 2
		case t.IsClosed():
			return 3
		default:
			return 4
		}
	}
	if c := cmp.Compare(pri(t), pri(o)); c != 0 {
		return c
	}
	return t.Type.Compare(o.Type)
}

func (s *Item) ActiveLinks() []Link {
	links := make([]Link, 0, len(s.Links))
	for _, link := range s.Links {
		if !link.IsDeleted() {
			links = append(links, link)
		}
	}
	return links
}

func (s *Item) PrintedLinks() []Link {
	links := make([]Link, 0, len(s.Links))
	for _, link := range s.Links {
		if !link.IsDeleted() && link.Public && strings.TrimSpace(link.Address) != "" {
			links = append(links, link)
		}
	}
	return links
}

func (s *Item) CreateItemLink() Link {
	var maxID int
	for i := range s.Links {
		link := s.Links[i]
		if link.ID > maxID {
			maxID = link.ID
		}
	}
	link := Link{ID: maxID + 1, Public: true}
	s.Links = append(s.Links, link)
	return link
}

func (s *Item) UpdateItemLink(link Link) {
	for i := range s.Links {
		if s.Links[i].ID == link.ID {
			s.Links[i] = link
			break
		}
	}
}

func (s *Item) DeleteItemLink(link Link, at time.Time) {
	for i := range s.Links {
		if s.Links[i].ID == link.ID {
			s.Links[i].DeletedAt = at
			break
		}
	}
}

func (s *Item) TitleForView() string {
	if len(s.Title) > 0 {
		return s.Title
	}
	return "no title"
}

func (s *Link) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}

func (s *VM) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}
