package graph

import (
	"context"
	"testing"

	"github.com/VadimRight/GraphQLOzon/internal/usecase"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/stretchr/testify/assert"
)

func TestPosts(t *testing.T) {
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{PostUsecase: mockPostUsecase, UserUsecase: mockUserUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedPosts := []*model.Post{
		{ID: "1", Text: "Test post", AuthorID: "1"},
	}

	expectedUser := &model.User{ID: "1", Username: "user1"}
	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Test comment", AuthorID: "1"},
	}

	mockPostUsecase.On("GetAllPosts", ctx, &limit, &offset).Return(expectedPosts, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByPostID", ctx, "1", &limit, &offset).Return(expectedComments, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByParentID", ctx, "1", &limit, &offset).Return([]*model.CommentResponse{}, nil)

	posts, err := resolver.Posts(ctx, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	mockPostUsecase.AssertExpectations(t)
	mockUserUsecase.AssertExpectations(t)
	mockCommentUsecase.AssertExpectations(t)
}

func TestPostsByUserID(t *testing.T) {
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{PostUsecase: mockPostUsecase, UserUsecase: mockUserUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	userID := "1"
	limit := 10
	offset := 0

	expectedPosts := []*model.Post{
		{ID: "1", Text: "Test post", AuthorID: "1"},
	}

	expectedUser := &model.User{ID: "1", Username: "user1"}
	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Test comment", AuthorID: "1"},
	}

	mockPostUsecase.On("GetPostsByUserID", ctx, userID, &limit, &offset).Return(expectedPosts, nil)
	mockUserUsecase.On("GetUserByID", ctx, userID).Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByPostID", ctx, "1", &limit, &offset).Return(expectedComments, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByParentID", ctx, "1", &limit, &offset).Return([]*model.CommentResponse{}, nil)

	posts, err := resolver.PostsByUserID(ctx, userID, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	mockPostUsecase.AssertExpectations(t)
	mockUserUsecase.AssertExpectations(t)
	mockCommentUsecase.AssertExpectations(t)
}

func TestPost(t *testing.T) {
	mockPostUsecase := new(usecase.MockPostUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	resolver := &queryResolver{&Resolver{PostUsecase: mockPostUsecase, UserUsecase: mockUserUsecase, CommentUsecase: mockCommentUsecase}}

	ctx := context.Background()
	id := "1"
	limit := 10
	offset := 0

	expectedPost := &model.Post{ID: "1", Text: "Test post", AuthorID: "1"}
	expectedUser := &model.User{ID: "1", Username: "user1"}
	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Test comment", AuthorID: "1"},
	}

	mockPostUsecase.On("GetPostByID", ctx, id).Return(expectedPost, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByPostID", ctx, "1", &limit, &offset).Return(expectedComments, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(expectedUser, nil)
	mockCommentUsecase.On("GetCommentsByParentID", ctx, "1", &limit, &offset).Return([]*model.CommentResponse{}, nil)

	post, err := resolver.Post(ctx, id, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	mockPostUsecase.AssertExpectations(t)
	mockUserUsecase.AssertExpectations(t)
	mockCommentUsecase.AssertExpectations(t)
}
