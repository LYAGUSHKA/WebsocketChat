package sqlstorage

import (
	"log"

	"github.com/Garius6/websocket_chat/model"
)

type MessageRepository struct {
	Storage *Storage
}

func (r *MessageRepository) Create(m *model.Message) error {
	log.Println(*m)
	_, err := r.Storage.Db.Exec(
		"INSERT INTO messages(m_from, data, m_to, chat_id) VALUES(?,?,?,?)",
		m.From,
		m.Data,
	)

	if err != nil {
		return nil
	}

	return nil
}

//GetLastMessages return last count messages from chatID chat
func (r *MessageRepository) GetLastMessagesInChat(chatID, offset, count int) ([]model.Message, error) {
	msgs := make([]model.Message, 0)
	messages, err := r.Storage.Db.Query(
		"SELECT * FROM messages WHERE chat_id=$1 LIMIT $2, $3",
		chatID,
		offset,
		count,
	)
	if err != nil {
		return nil, err
	}

	for messages.Next() {
		m := model.Message{}
		err = messages.Scan(&m.From, &m.Data)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}
