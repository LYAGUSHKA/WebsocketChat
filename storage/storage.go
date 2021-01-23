package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //...
)

type Storage struct {
	config         *Config
	Db             *sql.DB
	userRepository *UserRepository
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

//Message ...
type Message struct {
	ID      int
	Message string
}

func (s *Storage) SaveMessage(msg []byte) {
	message := string(msg)
	_, err := s.Db.Exec("INSERT INTO message(message) VALUES($1)", message)
	if err != nil {
		_ = fmt.Errorf("Storage: %s", err)
	}
}

func (s *Storage) GetLastMessages(number int) []Message {

	rows, err := s.Db.Query("SELECT * FROM message ORDER BY id ASC LIMIT 2, ?", number)
	if err != nil {
		_ = fmt.Errorf("Storage: %s", err)
	}
	var messages []Message

	for rows.Next() {
		s := Message{}
		err = rows.Scan(&s.ID, &s.Message)
		if err != nil {
			_ = fmt.Errorf("Rows: %s", err)
			continue
		}
		messages = append(messages, s)
	}

	return messages
}
