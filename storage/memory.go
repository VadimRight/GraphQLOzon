package storage

import (
	"context"
	"fmt"
	"sync"
	"github.com/VadimRight/GraphQLOzon/graph/model"
)

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

func (s *InMemoryStorage) CreateUser(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryStorage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryStorage) GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var posts []*model.Post
	for _, post := range s.posts {
		if post.AuthorID == userID {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (s *InMemoryStorage) CreateComment(ctx context.Context, comment *model.CommentResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[comment.ID] = comment
	return nil
}

func (s *InMemoryStorage) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*model.CommentResponse
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (s *InMemoryStorage) GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*model.CommentResponse
	for _, comment := range s.comments {
		if comment.ParentCommentID != nil && *comment.ParentCommentID == parentID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}
