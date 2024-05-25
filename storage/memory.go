package storage

import (
	"context"
	"fmt"
	"sync"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/google/uuid"
	"errors"
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

// Реализация методов интерфейса Storage для in-memory
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

func (s *InMemoryStorage) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]*model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *InMemoryStorage) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var posts []*model.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	// Пагинация
	if limit != nil && offset != nil {
		start := *offset
		end := *offset + *limit
		if start > len(posts) {
			return []*model.Post{}, nil
		}
		if end > len(posts) {
			end = len(posts)
		}
		posts = posts[start:end]
	}

	return posts, nil
}

func (s *InMemoryStorage) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var posts []*model.Post
	for _, post := range s.posts {
		if post.AuthorID == userID {
			posts = append(posts, post)
		}
	}

	// Пагинация
	if limit != nil && offset != nil {
		start := *offset
		end := *offset + *limit
		if start > len(posts) {
			return []*model.Post{}, nil
		}
		if end > len(posts) {
			end = len(posts)
		}
		posts = posts[start:end]
	}

	return posts, nil
}

func (s *InMemoryStorage) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, exists := s.posts[postID]
	if !exists {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	post := &model.Post{ID: id, Text: text, AuthorID: authorID, Commentable: commentable}
	s.posts[id] = post
	return post, nil
}

func (s *InMemoryStorage) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	comments := make([]*model.CommentResponse, 0, len(s.comments))
	for _, comment := range s.comments {
		comments = append(comments, comment)
	}
	return comments, nil
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

func (s *InMemoryStorage) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*model.CommentResponse
	for _, comment := range s.comments {
		if comment.ParentCommentID != nil && *comment.ParentCommentID == parentID {
			comments = append(comments, comment)
		}
	}

	// Применяем лимит и смещение, если они заданы
	if limit != nil && offset != nil {
		start := *offset
		end := start + *limit
		if start > len(comments) {
			return []*model.CommentResponse{}, nil
		}
		if end > len(comments) {
			end = len(comments)
		}
		return comments[start:end], nil
	}

	return comments, nil
}

func (s *InMemoryStorage) GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	comment, exists := s.comments[id]
	if !exists {
		return nil, fmt.Errorf("comment not found")
	}
	return comment, nil
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

func (s *InMemoryStorage) CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var isReply bool
	var parentCommentID *string
	var postID string
	var commentAble bool

	// Сначала проверяем, является ли itemId постом и включены ли комментарии
	if post, exists := s.posts[itemId]; exists {
		// itemId является постом
		postID = itemId
		isReply = false
		commentAble = post.Commentable
		if !commentAble {
			return nil, errors.New("author turned off comments under this post")
		}
	} else if comment, exists := s.comments[itemId]; exists {
		// itemId является комментарием
		postID = comment.PostID
		parentCommentID = &itemId
		isReply = true
	} else {
		return nil, errors.New("item not found")
	}

	id := uuid.New().String()
	var newComment *model.CommentResponse
	if isReply {
		// Если это ответ на комментарий
		newComment = &model.CommentResponse{ID: id, Comment: commentText, AuthorID: userID, PostID: postID, ParentCommentID: parentCommentID}
	} else {
		// Если это комментарий к посту
		newComment = &model.CommentResponse{ID: id, Comment: commentText, AuthorID: userID, PostID: postID}
	}
	s.comments[id] = newComment

	return newComment, nil
}
