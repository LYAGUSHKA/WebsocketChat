package storage

import "github.com/Garius6/websocket_chat/model"

type UserRepository struct {
	Storage *Storage
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	result, err := r.Storage.Db.Exec(
		"INSERT INTO users(nickname, encrypted_password) VALUES($1, $2); SELECT last_insert_rowid()",
		u.Nickname,
		u.EncryptedPassword,
	)
	if err != nil {
		return nil, err
	}
	buff, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = uint64(buff)
	return u, nil
}

func (r *UserRepository) FindByNickname(nickname string) (*model.User, error) {
	u := &model.User{}
	if err := r.Storage.Db.QueryRow(
		"SELECT * FROM users WHERE nickname = $1",
		nickname,
	).Scan(
		&u.ID,
		&u.Nickname,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}
	return u, nil
}
