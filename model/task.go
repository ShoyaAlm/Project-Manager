package model

type Category string

const (
	Strategic   Category = "strategic"
	Operational Category = "operational"
	Compliance  Category = "compliance"
)

type Priority int

const (
	Priority1 Priority = 1
	Priority2 Priority = 2
	Priority3 Priority = 3
	Priority4 Priority = 4
	Priority5 Priority = 5
)

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Category  Category
	Priority  Priority
}
