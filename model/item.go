package model

type Item struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	DueDate    string    `json:"duedate"`
	AssignedTo []*Member `json:"assignedto"`
	Done       bool      `json:"done"`
}
