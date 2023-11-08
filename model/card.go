package model

import "time"

// type Date time.Time

type Card struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Dates       []time.Time  `json:"dates"`
	Checklists  []*Checklist `json:"checklists"`
	Members     []*Member    `json:"members"`
	Owner       *User        `json:"owner"`
}
