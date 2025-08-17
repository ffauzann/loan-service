package model

type Message struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}
