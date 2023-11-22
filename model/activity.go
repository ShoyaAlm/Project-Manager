package model

import "time"

type Activity struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	CardID    int       `json:"card_id"`
	CreatedAt time.Time `json:"created_at"`
}

