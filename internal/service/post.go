package service

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
)

type PostService interface {
	CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error)
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
}

type postService struct {
	storage storage.Storage
}

func NewPostService(storage storage.Storage) PostService {
	return &postService{storage: storage}
}

func (s *postService) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	return s.storage.CreatePost(ctx, id, text, authorID, commentable)
}

func (s *postService) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	return s.storage.GetPostByID(ctx, id)
}

func (s *postService) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	return s.storage.GetPostsByUserID(ctx, userID, limit, offset)
}

func (s *postService) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	return s.storage.GetAllPosts(ctx, limit, offset)
}
