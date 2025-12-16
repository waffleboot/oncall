package model

import "time"

type Note struct {
	Text      string
	DeletedAt time.Time
}
