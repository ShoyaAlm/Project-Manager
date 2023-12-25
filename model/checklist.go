package model

type Checklist struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Items    []*Item `json:"items"`
	Position int     `json:"position"`
}
