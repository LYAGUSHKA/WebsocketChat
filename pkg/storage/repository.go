package storage

import "github.com/Garius6/websocket_chat/model"

type UserRepository interface {
	Create(login, password string) (*model.User, error)
	FindByID(ID int) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
}

type MessageRepository interface {
	Create(*model.Message) error
	GetLastMessagesInChat(chatID, offset, count int) ([]model.Message, error)
}
