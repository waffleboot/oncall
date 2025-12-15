package model

import (
	"cmp"
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
	VersionedObj[T any] struct {
		versions []T
	}
	ItemType string
	Item     struct {
		ID       uuid.UUID
		Num      int
		SleepAt  time.Time
		ClosedAt time.Time
		Type     ItemType
		Links    []ItemLink
	}
	ItemLink struct {
		ID          int
		Link        string
		Public      bool
		DeletedAt   time.Time
		Description VersionedObj[string]
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
		case ItemTypeInc:
			return 1
		case ItemTypeAsk:
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

func (s *Item) ActiveLinks() []ItemLink {
	links := make([]ItemLink, 0, len(s.Links))
	for i := range s.Links {
		if !s.Links[i].IsDeleted() {
			links = append(links, s.Links[i])
		}
	}
	return links
}

func (s *Item) CreateItemLink() ItemLink {
	var maxID int
	for i := range s.Links {
		link := s.Links[i]
		if link.ID > maxID {
			maxID = link.ID
		}
	}
	return ItemLink{ID: maxID + 1, Public: true}
}

func (s *ItemLink) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}

func NewVersionedObj[T any](versions []T) VersionedObj[T] {
	return VersionedObj[T]{versions: versions}
}

func (v *VersionedObj[T]) Value() T {
	var zero T
	if len(v.versions) == 0 {
		return zero
	}
	return v.versions[len(v.versions)-1]
}

func (v *VersionedObj[T]) SetValue(value T) {
	v.versions = append(v.versions, value)
}

func (v *VersionedObj[T]) Versions() []T {
	return v.versions
}
