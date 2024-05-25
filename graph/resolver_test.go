package graph

import (
	"context"
	"testing"

	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) HashPassword(password string) string {
	args := m.Called(password)
	return args.String(0)
}

func (m *MockUserService) ComparePassword(hashed string, normal string) error {
	args := m.Called(hashed, normal)
	return args.Error(0)
}

func (m *MockUserService) GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockUserService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockUserService) GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockUserService) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.User), args.Error(1)
}

func TestUsers(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()

	expectedUsers := []*model.User{
		{ID: "1", Username: "user1"},
		{ID: "2", Username: "user2"},
	}

	mockUserService.On("GetAllUsers", ctx).Return(expectedUsers, nil)

	users, err := resolver.Query().Users(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockUserService.AssertExpectations(t)
}

func TestUser(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()
	userID := "1"

	expectedUser := &model.User{ID: userID, Username: "user1"}

	mockUserService.On("GetUserByID", ctx, userID).Return(expectedUser, nil)

	user, err := resolver.Query().User(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserService.AssertExpectations(t)
}

func TestUserByUsername(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()
	username := "user1"

	expectedUser := &model.User{ID: "1", Username: username}

	mockUserService.On("GetUserByUsername", ctx, username).Return(expectedUser, nil)

	user, err := resolver.Query().UserByUsername(ctx, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserService.AssertExpectations(t)
}
