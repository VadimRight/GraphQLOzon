// storage/storage.go
package storage

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
)

type Storage interface {
	// Пользователи
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	// Посты
	GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postID string) (*model.Post, error)
	CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error)

	// Комментарии
	GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error)
	CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error)
}