package model

type Message struct {
	From   string `json:"from"`
	Data   string `json:"data"`
	To     string `json:"to"`
	ChatID int    `json:"chat_id"`
	//Timestamp string `json:"timestamp"`
}

func NewMessage() *Message {
	return &Message{}
}
