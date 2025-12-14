package model

const (
	ItemTypeInc   ItemType = "inc"
	ItemTypeAdhoc ItemType = "adhoc"
	ItemTypeAsk   ItemType = "ask"
	ItemTypeAlert ItemType = "alert"
)

type (
	ItemType string
	Item     struct {
		ID   int
		Type ItemType
	}
)

func (t ItemType) String() string {
	switch t {
	case ItemTypeInc:
		return "Инцидент"
	case ItemTypeAdhoc:
		return "Adhoc"
	case ItemTypeAsk:
		return "Обращение"
	case ItemTypeAlert:
		return "Alert"
	}
	return "Unknown"
}
