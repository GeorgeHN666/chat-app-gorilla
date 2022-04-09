package models

type Message struct {
	ID     int    `json:"id,omitempty"`
	Sender string `json:"sender"`
	Target string `json:"target"`
	Body   string `json:"body"`
}
