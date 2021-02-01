package model

type Message struct {
	From   string `json:"from"`
	Data   string `json:"data"`
	To     string `json:"to"`
	ChatID string `json:"chat_id"`
}

func NewMessage() *Message {
	return &Message{}
}
