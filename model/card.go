package model

import "time"

// type Date time.Time

type Card struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Dates       []time.Time  `json:"dates"`
	Checklists  []*Checklist `json:"checklists"`
	Members     []*User      `json:"members"`
	OwnerID     int          `json:"owner_id"`
	Owner       *User        `json:"owner"`
}
