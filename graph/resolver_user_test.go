package graph

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/VadimRight/GraphQLOzon/internal/usecase"
	"github.com/VadimRight/GraphQLOzon/model"
)

func TestUsers(t *testing.T) {
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{UserUsecase: mockUserUsecase, PostUsecase: mockPostUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedUsers := []*model.User{
		{ID: "1", Username: "user1"},
	}

	mockUserUsecase.On("GetAllUsers", ctx).Return(expectedUsers, nil)
	mockPostUsecase.On("GetPostsByUserID", ctx, "1", &limit, &offset).Return([]*model.Post{}, nil)

	users, err := resolver.Users(ctx, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockUserUsecase.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestUser(t *testing.T) {
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{UserUsecase: mockUserUsecase, PostUsecase: mockPostUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	id := "1"
	limit := 10
	offset := 0

	expectedUser := &model.User{ID: "1", Username: "user1"}

	mockUserUsecase.On("GetUserByID", ctx, id).Return(expectedUser, nil)
	mockPostUsecase.On("GetPostsByUserID", ctx, "1", &limit, &offset).Return([]*model.Post{}, nil)

	user, err := resolver.User(ctx, id, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserUsecase.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestUserByUsername(t *testing.T) {
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{UserUsecase: mockUserUsecase, PostUsecase: mockPostUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	username := "user1"
	limit := 10
	offset := 0

	expectedUser := &model.User{ID: "1", Username: "user1"}

	mockUserUsecase.On("GetUserByUsername", ctx, username).Return(expectedUser, nil)
	mockPostUsecase.On("GetPostsByUserID", ctx, "1", &limit, &offset).Return([]*model.Post{}, nil)

	user, err := resolver.UserByUsername(ctx, username, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserUsecase.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestLoginUser(t *testing.T) {
	mockUserUsecase := new(usecase.MockUserUsecase)
	resolver := &mutationResolver{&Resolver{UserUsecase: mockUserUsecase}}

	ctx := context.Background()
	username := "user1"
	password := "password"
	hashedPassword := "hashedPassword"
	userID := "1"

	expectedUser := &model.User{ID: "1", Username: "user1", Password: hashedPassword}
	expectedToken := "token"

	mockUserUsecase.On("GetUserByUsername", ctx, username).Return(expectedUser, nil)
	mockUserUsecase.On("ComparePassword", hashedPassword, password).Return(false)
	mockUserUsecase.On("GenerateToken", ctx, userID).Return(expectedToken, nil)

	token, err := resolver.LoginUser(ctx, username, password)

	assert.NoError(t, err)
	assert.Equal(t, &model.Token{Token: expectedToken}, token)
	mockUserUsecase.AssertExpectations(t)
}

