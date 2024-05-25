// storage/storage.go
package storage

import (
    	"context"
    	"github.com/VadimRight/GraphQLOzon/graph/model"
)

type Storage interface {
    	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
    	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
    	GetUserByID(ctx context.Context, userID string) (*model.User, error)
    	GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error)
    	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error)
    	GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error)
    	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
    	GetPostIdByItemId(ctx context.Context, itemId string) (*model.Post, error)
}
