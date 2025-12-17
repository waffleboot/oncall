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

func (s *Link) ToPublish() string {
	description := strings.TrimSpace(s.Description)
	if description != "" {
		return fmt.Sprintf("%s - %s", s.Address, description)
	}
	return s.Address
}

func (s *Link) MenuItem() string {
	var sb strings.Builder
	if s.Address == "" {
		sb.WriteString("empty")
	} else {
		sb.WriteString(s.Address)
	}
	if s.Public {
		sb.WriteString(" - public")
	} else {
		sb.WriteString(" - private")
	}
	return sb.String()
}

func (s *Link) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}
