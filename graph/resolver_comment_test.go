package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentService is a mock implementation of the CommentService interface
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockCommentService) GetCommentByID(ctx context.Context, id string) (*model.CommentResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.CommentResponse), args.Error(1)
}

func (m *MockCommentService) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, parentID, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockCommentService) CreateComment(ctx context.Context, commentText, itemId, userID string) (*model.CommentResponse, error) {
	args := m.Called(ctx, commentText, itemId, userID)
	return args.Get(0).(*model.CommentResponse), args.Error(1)
}


func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

// Tests for CommentService methods

func TestComments(t *testing.T) {
	mockCommentService := new(MockCommentService)
	mockUserService := new(MockUserService)
	resolver := &Resolver{
		CommentService: mockCommentService,
		UserService:    mockUserService,
	}

	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Test comment 1", AuthorID: "user1"},
		{ID: "2", Comment: "Test comment 2", AuthorID: "user2"},
	}

	mockCommentService.On("GetAllComments", mock.Anything, mock.Anything, mock.Anything).Return(expectedComments, nil)
	mockUserService.On("GetUserByID", mock.Anything, "user1").Return(&model.User{ID: "user1", Username: "user1"}, nil)
	mockUserService.On("GetUserByID", mock.Anything, "user2").Return(&model.User{ID: "user2", Username: "user2"}, nil)
	mockCommentService.On("GetCommentsByParentID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.CommentResponse{}, nil)

	limit := 10
	offset := 0
	comments, err := resolver.Query().Comments(context.Background(), &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	mockCommentService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestComment(t *testing.T) {
	mockCommentService := new(MockCommentService)
	mockUserService := new(MockUserService)
	resolver := &Resolver{
		CommentService: mockCommentService,
		UserService:    mockUserService,
	}

	ctx := context.Background()
	commentID := "1"

	expectedComment := &model.CommentResponse{ID: commentID, Comment: "Test comment", AuthorID: "user1"}

	mockCommentService.On("GetCommentByID", ctx, commentID).Return(expectedComment, nil)
	mockUserService.On("GetUserByID", ctx, "user1").Return(&model.User{ID: "user1", Username: "user1"}, nil)
	mockCommentService.On("GetCommentsByParentID", ctx, commentID, mock.Anything, mock.Anything).Return([]*model.CommentResponse{}, nil)

	limit := 10
	offset := 0
	comment, err := resolver.Query().Comment(ctx, commentID, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	mockCommentService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestCreateComment(t *testing.T) {
	mockCommentService := new(MockCommentService)
	mockUserService := new(MockUserService)
	resolver := &Resolver{
		CommentService: mockCommentService,
		UserService:    mockUserService,
	}

	ctx := context.Background()
	commentText := "New comment"
	itemID := "post1"
	userID := "user1"

	expectedComment := &model.CommentResponse{ID: "1", Comment: commentText, AuthorID: userID}

	mockCommentService.On("CreateComment", ctx, commentText, itemID, userID).Return(expectedComment, nil)
	mockUserService.On("GetUserByID", ctx, userID).Return(&model.User{ID: userID, Username: "user1"}, nil)

	comment, err := resolver.Mutation().CreateComment(ctx, commentText, itemID)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	mockCommentService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
