package port

type OnCallService interface {
	AddItem()
	Items() []string
}
