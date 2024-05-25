// service/comment_service.go
package service

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
)

type CommentService interface {
	GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error)
	CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error)
}

type commentService struct {
	storage storage.Storage
}

func NewCommentService(storage storage.Storage) CommentService {
	return &commentService{storage: storage}
}

func (s *commentService) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	return s.storage.GetAllComments(ctx, limit, offset)
}

func (s *commentService) GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error) {
	return s.storage.GetCommentByID(ctx, id)
}

func (s *commentService) CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error) {
	return s.storage.CreateComment(ctx, commentText, itemId, userID)
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByPostID(ctx, postID)
}

func (s *commentService) GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByParentID(ctx, parentID)
}
