package model

import "time"

type Node struct {
	ID        int
	Name      string
	DeletedAt time.Time
}
