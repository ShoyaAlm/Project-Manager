package model

import "time"

type Card struct {
	ID          int          	`json:"id"`
	Name        string       	`json:"name"`
	Description string       	`json:"description"`
	Dates       []time.Time  	`json:"dates"`
	Checklists  []*Checklist 	`json:"checklists"`
	Members     []*User      	`json:"members"`
	OwnerID     int	          	`json:"owner_id"`
	Owner       *User        	`json:"owner"`
	Label       *string       	`json:"label"` // Added 'Label' attribute
	Activity    *[]string       `json:"activity"` // Added 'Activity' attribute
	Position 	int				`json:"position"`
}
