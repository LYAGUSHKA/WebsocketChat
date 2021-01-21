package main

import (
	"database/sql"

	"github.com/Garius6/websocket_chat/storage"
)

//Room ...
type Room struct {
	Clients map[*Client]bool
	events  chan Event
	db      *sql.DB
}

func newRoom(db *sql.DB) *Room {
	return &Room{
		Clients: make(map[*Client]bool),
		events:  make(chan Event),
		db:      db,
	}
}

func (h *Room) run() {
	for {
		e := <-h.events
		switch e.Type {
		case newMessage:
			for client := range h.Clients {
				client.send <- e.Data.([]byte)
				storage.SaveMessage(h.db, e.Data.([]byte))
			}
		case register:
			client := e.Data.(*Client)
			h.Clients[client] = true
			for _, msg := range storage.GetLastMessages(h.db, 5) {
				client.send <- []byte(msg.Message)
			}
		case unregister:
			if _, ok := h.Clients[e.Data.(*Client)]; ok {
				delete(h.Clients, e.Data.(*Client))
				close(e.Data.(*Client).send)
			}
		}

	}
}
