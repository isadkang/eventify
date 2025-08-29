package models

import "time"

type Ticket struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	EventID   int       `json:"event_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"date"`
}
