package sockets

import (
	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/storage"
	"log"
)

//Room ...
type Room struct {
	Clients map[*Client]bool
	Storage *storage.Storage
	Events  chan Event
}

func NewRoom(db *storage.Storage) *Room {
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
			for client := range h.Clients {
				m := e.Data.(model.Message)
				m2 := model.NewMessage()
				m2.Data = m.Data
				client.send <- m
				//h.Storage.Db.Exec("INSERT INTO messages(m_from, data, m_to, chat_id) VALUES(?,?,?,?)", m.From, m.Data, m.To, m.ChatID)
				log.Println(h.Storage.User().Storage == nil)
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
