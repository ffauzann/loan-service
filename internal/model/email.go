package model

type EmailRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
