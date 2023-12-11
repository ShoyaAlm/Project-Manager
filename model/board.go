package model

type Board struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Lists []*List `json:"lists"`
}
