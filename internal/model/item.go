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
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   time.Time
		ClosedAt    time.Time
		Type        ItemType
		Notes       []Note
		Nodes       []Node
		Links       []Link
		VMs         []VM
		Title       string
		Description string
	}
)

func NewItem(num int) Item {
	return Item{
		ID:        uuid.New(),
		Num:       num,
		Type:      ItemTypeAsk,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (i *Item) InProgress() bool {
	return !i.IsClosed() && i.SleepAt.IsZero()
}

func (i *Item) IsSleep() bool {
	return !i.IsClosed() && !i.SleepAt.IsZero()
}

func (i *Item) IsClosed() bool {
	return !i.ClosedAt.IsZero()
}

func (i *Item) IsDeleted() bool {
	return !i.DeletedAt.IsZero()
}

func (i *Item) Sleep() {
	i.SleepAt = time.Now()
}

func (i *Item) Awake() {
	i.SleepAt = time.Time{}
}

func (i *Item) Close() {
	i.ClosedAt = time.Now()
}

func (i *Item) Delete() {
	i.DeletedAt = time.Now()
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
		case t.InProgress():
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

func (s *Item) MenuItem() string {
	if len(s.Title) > 0 {
		return s.Title
	}
	return "no title"
}

func (s *Item) ToPrint() string {
	var sb strings.Builder

	sb.WriteString(s.MenuItem())
	// sb.WriteString(fmt.Sprintf(" #%d", s.Num))

	// switch {
	// case s.IsActive():
	// 	sb.WriteString(" (in progress)")
	// case s.IsSleep():
	// 	sb.WriteString(" (in progress, in waiting)")
	// }

	return sb.String()
}
