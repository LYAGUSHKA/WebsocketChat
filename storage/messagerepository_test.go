package storage_test

import (
	"testing"

	"github.com/Garius6/websocket_chat/model"
	"github.com/Garius6/websocket_chat/storage"
	"github.com/stretchr/testify/assert"
)

func TestMessageRepositoryCreate(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("messages")

	m, err := s.Message().Create(&model.Message{From: "Perdej", Data: "Hello", To: "me", ChatID: 1})
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestMessageRepositoryGetLastMessages(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("messages")

	for i := 0; i < 12; i++ {
		s.Message().Create(&model.Message{From: "Perdej", Data: "Hello", To: "me", ChatID: 1})
	}

	msgs, _ := s.Message().GetLastMessages(1, 0, 10)
	assert.Equal(t, 10, len(msgs))
}
