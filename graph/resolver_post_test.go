package graph

import (
	"context"
	"testing"

	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostService is a mock implementation of the PostService interface
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(ctx context.Context, id, text, authorID string, commentable bool) (*model.Post, error) {
	args := m.Called(ctx, id, text, authorID, commentable)
	post := args.Get(0)
	if post == nil {
		return nil, args.Error(1)
	}
	return post.(*model.Post), args.Error(1)
}

func (m *MockPostService) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	args := m.Called(ctx, id)
	post := args.Get(0)
	if post == nil {
		return nil, args.Error(1)
	}
	return post.(*model.Post), args.Error(1)
}

func (m *MockPostService) GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}


// MockCommentService is a mock implementation of the CommentService interface
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*model.CommentResponse), args.Error(1)
}

func TestGetAllPosts(t *testing.T) {
	mockPostService := new(MockPostService)
	resolver := &Resolver{PostService: mockPostService}

	ctx := context.Background()

	expectedPosts := []*model.Post{
		{ID: "1", Text: "Post 1", AuthorID: "1", Commentable: true},
		{ID: "2", Text: "Post 2", AuthorID: "2", Commentable: false},
	}

	mockPostService.On("GetAllPosts", ctx).Return(expectedPosts, nil)

	posts, err := resolver.Query().Posts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	mockPostService.AssertExpectations(t)
}

func TestGetPostByID(t *testing.T) {
	mockPostService := new(MockPostService)
	mockUserService := new(MockUserService)
	mockCommentService := new(MockCommentService)

	resolver := &Resolver{
		PostService:    mockPostService,
		UserService:    mockUserService,
	}

	ctx := context.Background()
	postID := "1"
	authorID := "1"

	expectedPost := &model.Post{ID: postID, Text: "Post 1", AuthorID: authorID, Commentable: true}
	expectedAuthor := &model.User{ID: authorID, Username: "author"}
	expectedComments := []*model.CommentResponse{
		{ID: "1", Comment: "Comment 1", AuthorID: authorID, PostID: postID},
	}

	mockPostService.On("GetPostByID", ctx, postID).Return(expectedPost, nil)
	mockUserService.On("GetUserByID", ctx, authorID).Return(expectedAuthor, nil)
	mockCommentService.On("GetCommentsByPostID", ctx, postID).Return(expectedComments, nil)

	post, err := resolver.Query().Post(ctx, postID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	assert.Equal(t, expectedAuthor, post.Author)
	assert.Equal(t, expectedComments, post.Comments)
	mockPostService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
	mockCommentService.AssertExpectations(t)
}

func TestGetPostsByUserID(t *testing.T) {
	mockPostService := new(MockPostService)
	resolver := &Resolver{PostService: mockPostService}

	ctx := context.Background()
	userID := "1"

	expectedPosts := []*model.Post{
		{ID: "1", Text: "Post 1", AuthorID: userID, Commentable: true},
		{ID: "2", Text: "Post 2", AuthorID: userID, Commentable: false},
	}

	mockPostService.On("GetPostsByUserID", ctx, userID).Return(expectedPosts, nil)

	posts, err := resolver.PostService.GetPostsByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	mockPostService.AssertExpectations(t)
}
