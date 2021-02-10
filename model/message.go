package model

type Message struct {
	From string `json:"from"`
	Data string `json:"data"`
}

func NewMessage() *Message {
	return &Message{}
}
