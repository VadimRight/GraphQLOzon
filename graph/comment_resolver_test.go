// resolver_comment_test.go
package graph

import (
	"context"
	"testing"

	"github.com/VadimRight/GraphQLOzon/internal/usecase"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/stretchr/testify/assert"
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
	mockUserUsecase.On("GetUserByID", ctx, "1").Return(&model.User{ID: "1", Username: "user1"}, nil)
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
