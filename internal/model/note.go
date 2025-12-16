package model

import "time"

type Note struct {
	ID        int
	Text      string
	DeletedAt time.Time
}
