package model

import "time"

type Node struct {
	Name      string
	DeletedAt time.Time
}
