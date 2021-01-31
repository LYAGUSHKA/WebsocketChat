package storage_test

import (
	"testing"

	"github.com/Garius6/websocket_chat/storage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryCreate(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create("Anatole", "123")
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
	_, err := s.User().FindByLogin(nickname)
	assert.Error(t, err)

	s.User().Create(
		nickname,
		"Pook",
	)
	u, err := s.User().FindByLogin(nickname)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
