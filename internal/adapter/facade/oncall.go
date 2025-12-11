package facade

type onCallService struct {
}

func NewOnCallService() *onCallService {
	return &onCallService{}
}

func (s *onCallService) Items() []string {
	return []string{"first", "second"}
}
