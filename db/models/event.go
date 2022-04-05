package models

type EventSource struct {
	Id        string      `json:"id"`
	Event     string      `json:"event"`
	CreatedAt uint64      `json:"created_at"`
	Service   string      `json:"service"`
	Payload   interface{} `json:"payload"`
}
