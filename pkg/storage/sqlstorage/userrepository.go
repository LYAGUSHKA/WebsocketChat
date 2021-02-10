package sqlstorage

import "github.com/Garius6/websocket_chat/model"

type UserRepository struct {
	Storage *Storage
}

func (r *UserRepository) Create(login string, password string) (*model.User, error) {
	u := &model.User{Login: login, Password: password}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	result, err := r.Storage.Db.Exec(
		"INSERT INTO users(nickname, encrypted_password) VALUES($1, $2); SELECT last_insert_rowid()",
		u.Login,
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

func (r *UserRepository) FindByLogin(nickname string) (*model.User, error) {
	u := &model.User{}
	if err := r.Storage.Db.QueryRow(
		"SELECT * FROM users WHERE nickname = $1",
		nickname,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(ID int) (*model.User, error) {
	u := &model.User{}
	if err := r.Storage.Db.QueryRow(
		"SELECT * FROM users WHERE id = $1",
		ID,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return u, nil
}
