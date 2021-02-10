package sockets

import (
	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/pkg/storage"
)

//Room ...
type Room struct {
	Clients map[*Client]bool
	Storage storage.Storage
	Events  chan Event
}

func NewRoom(db storage.Storage) *Room {
	return &Room{
		Clients: make(map[*Client]bool),
		Storage: db,
		Events:  make(chan Event),
	}
}

func (h *Room) RunRoom() {
	chats := make(map[int]*Room)
	for {
		e := <-h.Events
		switch e.Type {
		case NewMessage:
			m := e.Data.(model.Message)
			for client := range h.Clients {
				client.send <- m
			}
		case Register:
			client := e.Data.(*Client)
			h.Clients[client] = true
		case Unregister:
			if _, ok := h.Clients[e.Data.(*Client)]; ok {
				delete(h.Clients, e.Data.(*Client))
				close(e.Data.(*Client).send)
			}
		case NewChat:
			chatInfo := e.Data.(ChatInfo)
			chats[chatInfo.ID] = chatInfo.Chat
		}

	}
}
