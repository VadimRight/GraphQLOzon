package service

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
)

type UserService interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type userService struct {
	storage bootstrap.Storage
}

func NewUserService(storage bootstrap.Storage) UserService {
	return &userService{storage: storage}
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.storage.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
