package storage_test

import (
	"testing"

	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/storage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryCreate(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(&model.User{
		Nickname: "Anatole",
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)

	var nick string
	err = s.User().Storage.Db.QueryRow("SELECT encrypted_password FROM users WHERE nickname=$1", "Anatole").Scan(&nick)
	assert.NoError(t, err)

}

func TestUserRepositoryFindByID(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("users")

	nickname := "Alexxx"
	_, err := s.User().FindByNickname(nickname)
	assert.Error(t, err)

	s.User().Create(&model.User{
		Nickname:          nickname,
		EncryptedPassword: "Pook",
	})
	u, err := s.User().FindByNickname(nickname)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
