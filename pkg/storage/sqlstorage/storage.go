package sqlstorage

import (
	"database/sql"

	"github.com/Garius6/websocket_chat/pkg/storage"
)

type Storage struct {
	Db                *sql.DB
	userRepository    *UserRepository
	messageRepository *MessageRepository
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Db: db,
	}
}

func (s *Storage) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		Storage: s,
	}

	return s.userRepository
}

func (s *Storage) Message() storage.MessageRepository {
	if s.messageRepository != nil {
		return s.messageRepository
	}

	s.messageRepository = &MessageRepository{
		Storage: s,
	}

	return s.messageRepository
}
