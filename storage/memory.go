package storage

import (
	"context"
	"fmt"
	"sync"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/google/uuid"
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

func (s *InMemoryStorage) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    id := uuid.New().String()
    user := &model.User{ID: id, Username: username, Password: password}
    s.users[id] = user
    return user, nil
}

func (s *InMemoryStorage) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    user, exists := s.users[userID]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
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

func (s *InMemoryStorage) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var comments []*model.CommentResponse
    for _, comment := range s.comments {
        if comment.AuthorID == userID {
            comments = append(comments, comment)
        }
    }
    return comments, nil
}

func (s *InMemoryStorage) GetPostIdByItemId(ctx context.Context, itemId string) (*model.Post, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    comment, exists := s.comments[itemId]
    if !exists {
        return nil, fmt.Errorf("item not found")
    }
    post, exists := s.posts[comment.PostID]
    if !exists {
        return nil, fmt.Errorf("post not found")
    }
    return post, nil
}
