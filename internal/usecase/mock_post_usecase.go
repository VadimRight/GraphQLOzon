package usecase

import (
	"context"

	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/stretchr/testify/mock"
)

type MockPostUsecase struct {
	mock.Mock
}

func (m *MockPostUsecase) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	args := m.Called(ctx, id, text, authorID, commentable)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostUsecase) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostUsecase) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostUsecase) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}
