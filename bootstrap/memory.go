package bootstrap

import (
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"sync"
)

// Тип in-memory хранилища
type InMemoryStorage struct {
	users    map[string]*model.User
	posts    map[string]*model.Post
	comments map[string]*model.CommentResponse
	mu       sync.RWMutex
}

// Функция возвращающая объект InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:    make(map[string]*model.User),
		posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.CommentResponse),
	}
}

func InitInMemoryStorage() *InMemoryStorage {
	storage := NewInMemoryStorage()

	// Создание начальных данных, если нужно
	storage.mu.Lock()
	defer storage.mu.Unlock()
	return storage
}
