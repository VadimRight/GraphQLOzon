package usecase

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUsecase) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUsecase) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) ComparePassword(hashed string, normal string) bool {
	args := m.Called(hashed, normal)
	return args.Bool(0)
}

func (m *MockUserUsecase) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockUserUsecase) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUsecase) GenerateToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) ValidateToken(ctx context.Context, token string) (*jwt.Token, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(hashed string, normal string) error {
	args := m.Called(hashed, normal)
	return args.Error(0)
}
