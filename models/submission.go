package models

import "time"

type QuizSubmission struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	EventTitle  string    `json:"event_title"`
	EventID     int       `json:"event_id"`
	Score       *int      `json:"score,omitempty"`
	Status      string    `json:"status,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
}
