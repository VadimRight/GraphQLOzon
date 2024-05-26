// resolver_comment_test.go
package graph

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/VadimRight/GraphQLOzon/internal/usecase"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
)

func TestComments(t *testing.T) {
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	resolver := &queryResolver{&Resolver{CommentUsecase: mockCommentUsecase, UserUsecase: mockUserUsecase}}

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Test comment", AuthorID: "1"},
	}

	mockCommentUsecase.On("GetAllComments", ctx, &limit, &offset).Return(expectedComments, nil)
	mockUserUsecase.On("GetUserByID", ctx, "2").Return(&model.User{ID: "1", Username: "user1"}, nil)
	mockCommentUsecase.On("GetCommentsByParentID", ctx, "1", &limit, &offset).Return([]*model.CommentResponse{}, nil)

	comments, err := resolver.Comments(ctx, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	mockCommentUsecase.AssertExpectations(t)
	mockUserUsecase.AssertExpectations(t)
}

func TestComment(t *testing.T) {
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	resolver := &queryResolver{&Resolver{CommentUsecase: mockCommentUsecase, UserUsecase: mockUserUsecase}}

	ctx := context.Background()
	id := "1"
	limit := 10
	offset := 0

	expectedComment := &model.CommentResponse{ID: "1", Comment: "Test comment", AuthorID: "1"}

	mockCommentUsecase.On("GetCommentByID", ctx, id).Return(expectedComment, nil)
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(&model.User{ID: "1", Username: "user1"}, nil)
	mockCommentUsecase.On("GetCommentsByParentID", ctx, "1", &limit, &offset).Return([]*model.CommentResponse{}, nil)

	comment, err := resolver.Comment(ctx, id, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	mockCommentUsecase.AssertExpectations(t)
	mockUserUsecase.AssertExpectations(t)
}

func TestCreateComment(t *testing.T) {
	mockCommentUsecase := new(usecase.MockCommentUsecase)
	mockUserUsecase := new(usecase.MockUserUsecase)
	resolver := &mutationResolver{&Resolver{CommentUsecase: mockCommentUsecase, UserUsecase: mockUserUsecase}}

	ctx := context.Background()
	commentText := "New comment"
	itemId := "1"
	userID := "1"

	expectedComment := &model.CommentResponse{ID: "1", Comment: "New comment", AuthorID: userID, PostID: itemId}

	mockCommentUsecase.On("CreateComment", ctx, commentText, itemId, userID).Return(expectedComment, nil)

	comment, err := resolver.CreateComment(ctx, commentText, itemId)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	mockCommentUsecase.AssertExpectations(t)
}
