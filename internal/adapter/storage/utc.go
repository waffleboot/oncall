package storage

import "time"

func from(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	t = t.UTC()
	return &t
}

func to[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
