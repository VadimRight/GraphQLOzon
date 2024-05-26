package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (m *MockUserService) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	args := m.Called(ctx, username, password)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (m *MockUserService) HashPassword(password string) string {
	args := m.Called(password)
	return args.String(0)
}

func (m *MockUserService) ComparePassword(hashed string, normal string) error {
	args := m.Called(hashed, normal)
	return args.Error(0)
}

func (m *MockUserService) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockUserService) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

// MockPostService is a mock implementation of the PostService interface
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	args := m.Called(ctx, id, text, authorID, commentable)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostService) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostService) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostService) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}


func (m *MockCommentService) GetCommentsByPostID(ctx context.Context, postID string, limit, offset *int) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}


// Tests for UserService methods

func TestRegisterUser(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()
	username := "test"
	password := "test"

	expectedUser := &model.User{ID: "1", Username: username}

	// Настройка мока для метода GetUserByUsername
	mockUserService.On("GetUserByUsername", ctx, username).Return(nil, errors.New("user not found"))
	// Настройка мока для метода HashPassword
	mockUserService.On("HashPassword", password).Return("hashedpassword")
	// Настройка мока для метода UserCreate
	mockUserService.On("UserCreate", ctx, username, "hashedpassword").Return(expectedUser, nil)

	user, err := resolver.Mutation().RegisterUser(ctx, username, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserService.AssertExpectations(t)
}

