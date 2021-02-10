package storage

type Storage interface {
	User() UserRepository
	Message() MessageRepository
}
