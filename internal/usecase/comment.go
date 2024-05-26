package usecase

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
)

type CommentUsecase interface {
	GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error)
	CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error)
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error)
}

type commentUsecase struct {
	storage storage.Storage
}

func NewCommentUsecase(storage storage.Storage) CommentUsecase {
	return &commentUsecase{storage: storage}
}

func (s *commentUsecase) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	return s.storage.GetAllComments(ctx, limit, offset)
}

func (s *commentUsecase) GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error) {
	return s.storage.GetCommentByID(ctx, id)
}

func (s *commentUsecase) CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error) {
	return s.storage.CreateComment(ctx, commentText, itemId, userID)
}

func (s *commentUsecase) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByPostID(ctx, postID, limit, offset)
}

func (s *commentUsecase) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByParentID(ctx, parentID, limit, offset)
}
