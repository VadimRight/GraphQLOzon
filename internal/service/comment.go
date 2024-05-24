package service

import (
	
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/lib/pq"
)

//Интерфейс сервиса комментариев
type CommentService interface {
	GetPostIdByItemId(ctx context.Context, id string) (*model.CommentResponse, error)
}

// Тип сервиса комментариев для доступа к типу Storage
type commentService struct {
	storage bootstrap.Storage
}

// Функция инициализации сервиса комментариев для запуска сервера с GraphQL Playground, бизнес-логика которого храниться в graph/resolver.go, но вызов просиходит в bootstrap/api.go 
func NewCommentService(storage bootstrap.Storage) CommentService {
	return &commentService{storage: storage}
}

// Вспомогательная метод комментариев для проверки в файле graph/resolver.go в функции создания комментария - если она не вернет ошибку - то комментарий пишется для поста и запускается проверка разрешены ли коментарии для указанного поста
func (s *commentService) GetPostIdByItemId(ctx context.Context, id string) (*model.CommentResponse, error) {
	var comment model.CommentResponse
	err := s.storage.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", id).Scan(&comment.ID, &comment.Comment)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
    			return nil, postgresErr
		}				
		return nil, err
	}
	return &comment, nil
}
