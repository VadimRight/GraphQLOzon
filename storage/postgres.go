package storage

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"context"
	"github.com/google/uuid"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"log"
	"fmt"
)

// Тип базы данных
type PostgresStorage struct {
	DB *sql.DB
}


// Все SQL запросы и функции работы с базой данных храняться в файле graph/resolver.go, а также вспомогательные запросы для обеспечения функционала схемы и резольвера храняться в сервисе пользователей в internal/service/user.go

// Функция возвращающая объект PostgresStorage
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
    return &PostgresStorage{DB: db}
}

// Функция инициализации базы данных и подключение к базе данных
func InitPostgresDatabase(cfg *bootstrap.Config) *PostgresStorage  {
	const op = "postgres.InitPostgresDatabase"

	dbHost := cfg.Postgres.PostgresHost
	dbPort := cfg.Postgres.PostgresPort
	dbUser := cfg.Postgres.PostgresUser
	dbPasswd := cfg.Postgres.PostgresPassword
	dbName := cfg.Postgres.DatabaseName

	postgresUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",dbHost, dbPort, dbUser, dbPasswd, dbName)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	// Создание таблицы поользователя, у которго есть зашифрованный пароль, имя и уникальный ID 
	createUserTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(20) NOT NULL UNIQUE,
		password CHAR(60) NOT NULL UNIQUE
	);`)
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createUserTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	// Создание таблицы постов, у которых есть текст, уникальный ID, а также ID пользователя, написавшего пост
	createPostTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS post (
		id UUID PRIMARY KEY,
		text TEXT NOT NULL,
		author_id UUID NOT NULL,
		commentable BOOLEAN NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(id));
	`)	
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createPostTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	// Создание таблицы комментариев, у которых есть сам текст комменатария, ID пользователя, оставившего комментарий, а также есть ID поста, под которым комментарий был написан и это поле всегда заполяется даже если комментарий оставлен не прямо к посту, а также есть ID коментария - это поле заполяется только тогда, когда комментарий оставлен к другому коментарию. Такая конструкция сущности комментария позволяет нам создавать иерархическую структуру данных.
	createCommentTable, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comment (
		id UUID PRIMARY KEY,
		comment VARCHAR(2000),
		author_id UUID NOT NULL,
		post_id UUID NOT NULL,
		parent_comment_id UUID,
    		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES post(id),
		FOREIGN KEY (parent_comment_id) REFERENCES comment(id)
	);`)
	if err != nil {	log.Fatalf("%s: %v", op, err) }
	_, err = createCommentTable.Exec()
	if err != nil {	log.Fatalf("%s: %v", op, err) }

	return &PostgresStorage{DB: db}
}

// Функция закрытия соединения с базой данных
func (s * PostgresStorage) ClosePostgres(db *PostgresStorage) error {
	return db.DB.Close()
}

func (s *PostgresStorage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    err := s.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        if postgresErr, ok := err.(*pq.Error); ok {
            return nil, postgresErr
        }
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
