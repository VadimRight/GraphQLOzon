package usecase

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/stretchr/testify/mock"
)

type MockCommentUsecase struct {
	mock.Mock
}

func (m *MockCommentUsecase) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockCommentUsecase) GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.CommentResponse), args.Error(1)
}

func (m *MockCommentUsecase) CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error) {
	args := m.Called(ctx, commentText, itemId, userID)
	return args.Get(0).(*model.CommentResponse), args.Error(1)
}

func (m *MockCommentUsecase) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockCommentUsecase) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, parentID, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

