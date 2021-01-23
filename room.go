package main

import "github.com/Garius6/websocket_chat/storage"

const (
	NEWMESSAGE = iota
	unregister
	register
	NEWCHAT
)

//Event ...
type Event struct {
	Type int
	Data interface{}
	// Sender interface{}
}

//Room ...
type Room struct {
	Clients map[*Client]bool
	Storage *storage.Storage
	events  chan Event
}

func newRoom(db *storage.Storage) *Room {
	return &Room{
		Clients: make(map[*Client]bool),
		Storage: db,
		events:  make(chan Event),
	}
}

func (h *Room) run() {
	chats := make(map[int]*Room)
	for {
		e := <-h.events
		switch e.Type {
		case NEWMESSAGE:
			for client := range h.Clients {
				client.send <- e.Data.(Message)
				//storage.SaveMessage(h.db, e.Data.([]byte))
			}
		case register:
			client := e.Data.(*Client)
			h.Clients[client] = true
			// for _, msg := range storage.GetLastMessages(h.db, 5) {
			// 	client.send <- []byte(msg.Message)
			// }
		case unregister:
			if _, ok := h.Clients[e.Data.(*Client)]; ok {
				delete(h.Clients, e.Data.(*Client))
				close(e.Data.(*Client).send)
			}
		case NEWCHAT:
			chatInfo := e.Data.(struct {
				ID   int
				chat *Room
			})
			chats[chatInfo.ID] = chatInfo.chat
		}

	}
}
