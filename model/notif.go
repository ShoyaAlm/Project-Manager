package model

import "time"

type Notification struct {
	ID       	int       `json:"id"`
	Message  	string    `json:"message"`
	UserID   	int    	  `json:"user_id"`
	CreatedAt 	time.Time `json:"created_at"`
	Read     	bool      `json:"read"`
}
