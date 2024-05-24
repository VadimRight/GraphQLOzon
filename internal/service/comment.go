package service

import (
	
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/lib/pq"
)

type CommentService interface {
	GetPostIdByItemId(ctx context.Context, id string) (*model.CommentResponse, error)
}

type commentService struct {
	storage bootstrap.Storage
}

func NewCommentService(storage bootstrap.Storage) CommentService {
	return &commentService{storage: storage}
}

func (s *commentService) GetPostIdByItemId(ctx context.Context, id string) (*model.CommentResponse, error) {
	var comment model.CommentResponse
	err := s.storage.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", id).Scan(&comment.ID, &comment.Comment)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
    			return nil, postgresErr
		}				
		return nil, err
	}
	return &comment, nil
}
