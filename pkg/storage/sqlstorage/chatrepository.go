package sqlstorage

import "github.com/Garius6/websocket_chat/model"

type ChatRepository struct {
	storage *Storage
}

func (c *ChatRepository) FindByUserID(userID int) ([]*model.Chat, error) {
	rows, err := c.storage.Db.Query("SELECT chait_id FROM chat_user WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	chats := make([]*model.Chat, 0)
	for rows.Next() {
		c := &model.Chat{}
		err = rows.Scan(
			&c.ID,
		)
		if err != nil {
			continue
		}

		chats = append(chats, c)
	}

	return chats, nil
}
