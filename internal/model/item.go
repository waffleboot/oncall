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

var itemTypePriorities = map[ItemType]int{
	ItemTypeInc:   1,
	ItemTypeAsk:   2,
	ItemTypeAlert: 3,
	ItemTypeAdhoc: 4,
}

func (t ItemType) Compare(o ItemType) int {
	return cmp.Compare(itemTypePriorities[t], itemTypePriorities[o])
}
