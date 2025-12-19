package model

import (
	"errors"
	"time"
)

type Journal struct {
	Items []Item
	Next  *User
}

func NewJournal() Journal {
	return Journal{}
}

func (s *Journal) GetItem(id int) (Item, error) {
	if i, err := s.getItem(id); err != nil {
		return Item{}, err
	} else {
		return s.Items[i], nil
	}
}

func (s *Journal) CreateItem(num int) Item {
	item := NewItem(num)
	s.Items = append(s.Items, item)
	return item
}

func (s *Journal) UpdateItem(item Item) (Item, error) {
	if i, err := s.getItem(item.ID); err != nil {
		return Item{}, err
	} else {
		item.UpdatedAt = time.Now()
		s.Items[i] = item
		return item, nil
	}
}

func (s *Journal) getItem(id int) (int, error) {
	for i := range s.Items {
		if s.Items[i].ID == id {
			return i, nil
		}
	}
	return 0, errors.New("not found")
}
