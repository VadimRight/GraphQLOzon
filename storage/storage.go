package storage

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
)

type Storage interface {
	// Пользователи
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	// Посты
	GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error)
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postID string) (*model.Post, error)
	CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error)

	// Комментарии
	GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*model.CommentResponse, error) // Обновлено
	GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error)
	CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error)
}

func StorageType(cfg *bootstrap.Config) Storage {
	storageType := cfg.Storage.StorageType
	var storage Storage
	if storageType == "memory" {
		storage = InitInMemoryStorage()
	} else {
		storage = InitPostgresDatabase(cfg)
	}
	return storage
}
