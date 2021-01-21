package main

const (
	newMessage = iota
	unregister
	register
)

//Event ...
type Event struct {
	Type int
	Data interface{}
	// Sender interface{}
}
