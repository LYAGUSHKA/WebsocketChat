package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//Message ...
type Message struct {
	ID      int
	Message string
}

func SaveMessage(db *sql.DB, msg []byte) {
	message := string(msg)
	_, err := db.Exec("INSERT INTO message(message) VALUES($1)", message)
	if err != nil {
		_ = fmt.Errorf("Storage: %s", err)
	}
}

func GetLastMessages(db *sql.DB, number int) []Message {
	rows, err := db.Query("SELECT * FROM message ORDER BY id ASC LIMIT 2, ?", number)
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
