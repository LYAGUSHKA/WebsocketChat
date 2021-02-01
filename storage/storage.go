package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //...
)

type Storage struct {
	config            *Config
	Db                *sql.DB
	userRepository    *UserRepository
	messageRepository *MessageRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("sqlite3", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	s.Db = db
	return nil
}

func (s *Storage) Close() {
	s.Db.Close()
}

func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		Storage: s,
	}

	return s.userRepository
}

func (s *Storage) Message() *MessageRepository {
	if s.userRepository != nil {
		return s.messageRepository
	}

	s.messageRepository = &MessageRepository{
		Storage: s,
	}

	return s.messageRepository
}
