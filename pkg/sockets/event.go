package sockets

const (
	NewMessage = iota
	Unregister
	Register
	NewChat
)

//Event ...
type Event struct {
	Type int
	Data interface{}
	// Sender interface{}
}
