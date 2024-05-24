package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"github.com/lib/pq"
)

type UserService interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	HashPassword(password string) string
	ComparePassword(hashed string, normal string) error
}

type userService struct {
	storage bootstrap.Storage
}

func NewUserService(storage bootstrap.Storage) UserService {
	return &userService{storage: storage}
}

// Служебная функция для получения пользователя из базы данных при попытке авторизации и проверки наличия пользователя в базе данных
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.storage.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
    			return nil, postgresErr
		}				
		return nil, err
	}
	return &user, nil
}

func (s *userService) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
	id := uuid.New().String()
	_, err := s.storage.DB.ExecContext(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", id, username, password)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
    			return nil, postgresErr
		}				
		return nil, err
	}
	return &model.User{ID: id, Username: username, Password: password}, nil
}
