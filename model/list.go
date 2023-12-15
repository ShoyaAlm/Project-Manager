package model

type List struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Cards    []*Card `json:"cards"`
	Position int     `json:"position"`
}

// StartDate string  `json:"startDate"`
// EndDate   string  `json:"endDate"`
