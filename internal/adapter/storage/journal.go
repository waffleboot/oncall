package storage

type storedData struct {
	LastNum int    `json:"last_num,omitempty"`
	Items   []item `json:"items,omitempty"`
}
