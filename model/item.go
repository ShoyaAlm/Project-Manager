package model

import "time"

type Item struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	StartDate  time.Time `json:"startdate"`
	DueDate    time.Time `json:"duedate"`
	Done       bool      `json:"done"`
	AssignedTo []*Member  `json:"assignedto"`
	Position   int 		`json:"position"`
}
