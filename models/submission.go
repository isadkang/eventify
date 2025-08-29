package models

import "time"

type QuizSubmission struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	EventID     int        `json:"event_id"`
	Score       *int       `json:"score,omitempty"`
	Status      string     `json:"status,omitempty"`
	SubmittedAt time.Time  `json:"submitted_at"`
}
