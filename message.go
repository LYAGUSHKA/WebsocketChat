package main

type Message struct {
	Data string `json:"data"`
}

func newMessage() *Message {
	return &Message{
		Data: "",
	}
}
