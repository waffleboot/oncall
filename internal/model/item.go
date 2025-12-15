package model

import (
	"cmp"
	"time"
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
		ID       int
		SleepAt  time.Time
		ClosedAt time.Time
		Type     ItemType
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
