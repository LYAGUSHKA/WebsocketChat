package storage

import (
	"github.com/Garius6/websocket_chat/model"
)

type MessageRepository struct {
	Storage *Storage
}

func (r *MessageRepository) Create(m *model.Message) (*model.Message, error) {
	_, err := r.Storage.Db.Exec(
		"INSERT INTO messages(m_from, data, m_to, chat_id) VALUES(?,?,?,?)",
		m.From,
		m.Data,
		m.To,
		m.ChatID,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}

//GetLastMessages return last count messages from chatID chat
func (r *MessageRepository) GetLastMessages(chatID, offset, count int) ([]model.Message, error) {
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
		err = messages.Scan(&m.From, &m.Data, &m.To, &m.ChatID)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}
