package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"context"
	"github.com/google/uuid"
)

// Тип базы данных
type PostgresStorage struct {
	DB *sql.DB
}


// Все SQL запросы и функции работы с базой данных храняться в файле graph/resolver.go, а также вспомогательные запросы для обеспечения функционала схемы и резольвера храняться в сервисе пользователей в internal/service/user.go

func NewPostgresStorage(DB *sql.DB) *PostgresStorage {
    return &PostgresStorage{DB: DB}
}

func (s *PostgresStorage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
        	return nil, err
	}
	return &user, nil
}

func (s *PostgresStorage) UserCreate(ctx context.Context, username string, password string) (*model.User, error) {
    	id := uuid.New().String()
    	_, err := s.DB.ExecContext(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", id, username, password)
    	if err != nil {
		return nil, err
    	}
    	return &model.User{ID: id, Username: username, Password: password}, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
    	var user model.User
    	err := s.DB.QueryRowContext(ctx, "SELECT id, username FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Username)
    	if err != nil {
    	    return nil, err
    	}
    	return &user, nil
}

func (s *PostgresStorage) GetPostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
    	rows, err := s.DB.QueryContext(ctx, "SELECT id, text, author_id, commentable FROM post WHERE author_id=$1", userID)
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
    	    posts = append(posts, &post)
    	}
    	return posts, nil
}

func (s *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.CommentResponse, error) {
    	rows, err := s.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE post_id=$1", postID)
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

func (s *PostgresStorage) GetCommentsByParentID(ctx context.Context, parentID string) ([]*model.CommentResponse, error) {
    	rows, err := s.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE parent_comment_id=$1", parentID)
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

func (s *PostgresStorage) GetCommentsByUserID(ctx context.Context, userID string) ([]*model.CommentResponse, error) {
    	rows, err := s.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE author_id=$1", userID)
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

func (s *PostgresStorage) GetPostIdByItemId(ctx context.Context, itemId string) (*model.Post, error) {
    	var postID string
    	err := s.DB.QueryRowContext(ctx, "SELECT post_id FROM comment WHERE id=$1", itemId).Scan(&postID)
    	if err != nil {
    	    return nil, err
    	}
    	var post model.Post
    	err = s.DB.QueryRowContext(ctx, "SELECT id, text, author_id, commentable FROM post WHERE id=$1", postID).Scan(&post.ID, &post.Text, &post.AuthorID, &post.Commentable)
    	if err != nil {
    	    return nil, err
    	}
    	return &post, nil
}
