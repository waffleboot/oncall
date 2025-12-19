package storage

import "github.com/waffleboot/oncall/internal/model"

type journal struct {
	LastNum int    `json:"last_num,omitempty"`
	Items   []item `json:"items,omitempty"`
	Next    string `json:"next,omitempty"`
}

func (j *journal) toDomain(user *model.User) model.Journal {
	items := make([]model.Item, 0, len(j.Items))
	for _, item := range j.Items {
		items = append(items, item.toDomain())
	}
	return model.Journal{Items: items, Next: user}
}

func (j *journal) fromDomain(journal model.Journal) {
	if journal.Next != nil {
		j.Next = journal.Next.Nick
	}
	j.Items = make([]item, len(journal.Items))
	for i, item := range journal.Items {
		j.Items[i].fromDomain(item)
	}
}
