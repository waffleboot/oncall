package model

import (
	"fmt"
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

	fmt.Fprintf(sb, "#%d ", s.ID)

	if s.Address == "" {
		sb.WriteString("empty")
	} else {
		sb.WriteString(s.Address)
	}

	return sb.String()
}

func (s *Link) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}
