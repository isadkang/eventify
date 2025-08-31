package models

import "time"

type Ticket struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	EventID   int       `json:"event_id"`
	EventName string    `json:"event_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"date"`
}
