package usecase

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
	"github.com/VadimRight/GraphQLOzon/internal/service"
)

type UserUsecase interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashed string, normal string) error
	GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type userUsecase struct {
	storage        storage.Storage
	commentUsecase CommentUsecase
	passwordService service.PasswordService
	authService    service.AuthService
}

func NewUserUsecase(storage storage.Storage, commentUsecase CommentUsecase, passwordService service.PasswordService, authService service.AuthService) UserUsecase {
	return &userUsecase{
		storage:        storage,
		commentUsecase: commentUsecase,
		passwordService: passwordService,
		authService:    authService,
	}
}

func (s *userUsecase) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.storage.GetAllUsers(ctx)
}

func (s *userUsecase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}

func (s *userUsecase) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	hashedPassword, err := s.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.storage.UserCreate(ctx, username, hashedPassword)
}

func (s *userUsecase) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
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

func (s *userUsecase) GetPostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	posts, err := s.storage.GetPostsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		post.Comments, err = s.commentUsecase.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (s *userUsecase) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	return s.storage.GetCommentsByUserID(ctx, userID)
}

func (s *userUsecase) HashPassword(password string) (string, error) {
	return s.passwordService.HashPassword(password)
}

func (s *userUsecase) ComparePassword(hashed string, normal string) error {
	return s.passwordService.ComparePassword(hashed, normal)
}
