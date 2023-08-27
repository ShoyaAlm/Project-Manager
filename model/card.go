package model

type Card struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Dates       []string     `json:"dates"`
	Checklists  []*Checklist `json:"checklists"`
	Members     []*Member    `json:"members"`
}
