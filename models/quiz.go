package models

import "encoding/json"

type Quiz struct {
	ID        int    `json:"id"`
	EventID   int    `json:"event_id"`
	Question  string `json:"question"`
	Options   json.RawMessage `json:"options"`
	AnswerKey string `json:"answer_key"`
}
