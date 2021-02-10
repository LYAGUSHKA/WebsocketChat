package model

type Chat struct {
	ID       uint64
	Messages []Message
	Users    []User
}
