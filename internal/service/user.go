package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/storage"
	"github.com/lib/pq"
)

//Интерфейс сервиса пользователей
type UserService interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UserCreate(ctx context.Context, username string, password string) (*model.User, error)
	HashPassword(password string) string
	ComparePassword(hashed string, normal string) error
	GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error)
	GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

// Тип сервиса пользователей
type userService struct {
	storage storage.PostgresStorage
}

// Функция инициализации сервиса пользователей для запуска сервера с GraphQL Playground, бизнес-логика которого храниться в graph/resolver.go, но вызов просиходит в bootstrap/api.go 
func NewUserService(storage storage.PostgresStorage) UserService {
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

// Метод создания пользователя 
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

// Метод получения пользователя по его ID получения 
func (s *userService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := s.storage.DB.QueryRowContext(ctx, "SELECT id, username FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	user.Posts, err = s.GetPostsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Comments, err = s.GetCommentsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Метод получения поста по ID пользователя
func (s *userService) GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
	rows, err := s.storage.DB.QueryContext(ctx, "SELECT id, text, author_id, commentable FROM post WHERE author_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Text, &post.AuthorID, &post.Commentable); err != nil {
			return nil, err
		}

		post.Comments, err = s.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}
	return posts, nil
}

// Метод получения комментариев по ID поста
func (s *userService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
	rows, err := s.storage.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE post_id=$1", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.CommentResponse
	for rows.Next() {
		var comment model.CommentResponse
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.PostID, &comment.ParentCommentID); err != nil {
			return nil, err
		}

		comment.Replies, err = s.GetCommentsByParentID(ctx, comment.ID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}
	return comments, nil
}

// Метод получения комментария по ID родительского комментария (для комментариев оставленных для других комментариев)
func (s *userService) GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error) {
	rows, err := s.storage.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE parent_comment_id=$1", parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.CommentResponse
	for rows.Next() {
		var comment model.CommentResponse
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.PostID, &comment.ParentCommentID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

// Метод получения комментариев по ID пользователей
func (s *userService) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
	rows, err := s.storage.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE author_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.CommentResponse
	for rows.Next() {
		var comment model.CommentResponse
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.PostID, &comment.ParentCommentID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
