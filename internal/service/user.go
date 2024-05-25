package service

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	HashPassword(password string) string
	ComparePassword(hashed string, normal string) error
	GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type userService struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) UserService {
	return &userService{storage: storage}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.storage.GetAllUsers(ctx)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}

func (s *userService) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	return s.storage.UserCreate(ctx, username, password)
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Posts, err = s.GetPostsByUserID(ctx, user.ID, nil, nil)
	if err != nil {
		return nil, err
	}

	user.Comments, err = s.GetCommentsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	posts, err := s.storage.GetPostsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		post.Comments, err = s.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (s *userService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	comments, err := s.storage.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		comment.Replies, err = s.GetCommentsByParentID(ctx, comment.ID, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	return comments, nil
}

func (s *userService) GetCommentsByParentID(ctx context.Context, parentID string, limit, offset *int) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByParentID(ctx, parentID, limit, offset)
}

func (s *userService) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByUserID(ctx, userID)
}
