package facade

import "fmt"

type onCallService struct {
	items []string
}

func NewOnCallService() *onCallService {
	return &onCallService{}
}

func (s *onCallService) Items() []string {
	return s.items
}

func (s *onCallService) AddItem() {
	s.items = append(s.items, fmt.Sprintf("item %d", len(s.items)+1))
}
