package model

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Bio      string  `json:"bio"`
	Cards    []*Card `json:"cards"`
}
