package model

type Project struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	Tasks     []*Task `json:"tasks"`
	Members   []*User `json:"members"`
}
