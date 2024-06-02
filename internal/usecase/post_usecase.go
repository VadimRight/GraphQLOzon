package usecase

import (
	"context"

	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/VadimRight/GraphQLOzon/storage"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error)
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
}

type postUsecase struct {
	storage storage.Storage
}

func NewPostUsecase(storage storage.Storage) PostUsecase {
	return &postUsecase{storage: storage}
}

func (s *postUsecase) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	return s.storage.CreatePost(ctx, id, text, authorID, commentable)
}

func (s *postUsecase) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	return s.storage.GetPostByID(ctx, id)
}

func (s *postUsecase) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	return s.storage.GetPostsByUserID(ctx, userID, limit, offset)
}

func (s *postUsecase) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	return s.storage.GetAllPosts(ctx, limit, offset)
}
