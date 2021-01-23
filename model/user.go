package model

type User struct {
	ID                uint64
	Nickname          string
	EncryptedPassword string
}
