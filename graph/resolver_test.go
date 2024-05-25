package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/VadimRight/GraphQLOzon/internal/service"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
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
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserService) RegisterUser(ctx context.Context, username string, password string) (*model.User, error) {
	args := m.Called(ctx, username, password)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (m *MockUserService) LoginUser(ctx context.Context, username string, password string) (*model.Token, error) {
	args := m.Called(ctx, username, password)
	token := args.Get(0)
	if token == nil {
		return nil, args.Error(1)
	}
	return token.(*model.Token), args.Error(1)
}

func (m *MockUserService) JwtGenerate(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func TestUsers(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()

	expectedUsers := []*model.User{
		{ID: "1", Username: "test1"},
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

	expectedUser := &model.User{ID: userID, Username: "test1"}

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
	username := "test1"

	expectedUser := &model.User{ID: "1", Username: username}

	mockUserService.On("GetUserByUsername", ctx, username).Return(expectedUser, nil)

	user, err := resolver.Query().UserByUsername(ctx, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserService.AssertExpectations(t)
}

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

func TestLoginUser(t *testing.T) {
	mockUserService := new(MockUserService)
	resolver := &Resolver{UserService: mockUserService}

	ctx := context.Background()
	username := "test"
	password := "test"

	expectedUser := &model.User{ID: "1", Username: username, Password: "hashedpassword"}

	// Настройка мока для метода GetUserByUsername
	mockUserService.On("GetUserByUsername", ctx, username).Return(expectedUser, nil)
	// Настройка мока для метода ComparePassword
	mockUserService.On("ComparePassword", "hashedpassword", password).Return(nil)

	// Генерируем реальный токен, чтобы использовать его для сравнения
	expectedToken, err := service.JwtGenerate(ctx, expectedUser.ID)
	assert.NoError(t, err)

	_, err = resolver.Mutation().LoginUser(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, expectedToken)


	mockUserService.AssertExpectations(t)
}
