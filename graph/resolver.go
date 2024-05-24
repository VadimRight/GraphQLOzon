// resolver/resolver.go
package graph

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/google/uuid"
	"github.com/lib/pq"
)


type Resolver struct{
	UserService service.UserService
	DB *sql.DB
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Функция получения всех пользователей
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}

		user.Posts, err = r.UserService.GetPostsByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		user.Comments, err = r.UserService.GetCommentsByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username FROM users WHERE id=$1", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	user.Posts, err = r.UserService.GetPostsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Comments, err = r.UserService.GetCommentsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Функция получения пользователя по ID
func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	user.Posts, err = r.UserService.GetPostsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Comments, err = r.UserService.GetCommentsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Функция логина пользователя
func (r *mutationResolver) LoginUser(ctx context.Context, username string, password string) (*model.Token, error) {
	getUser, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			return nil, postgresErr
		}
		return nil, err
	}
	if err := r.UserService.ComparePassword(getUser.Password, password); err != nil {
		return nil, err
	}
	token, err := service.JwtGenerate(ctx, getUser.ID)
	if err != nil {
		return nil, err
	}
	return &model.Token{Token: token}, nil
}

// Функция регистрации пользователя
func (r *mutationResolver) RegisterUser(ctx context.Context, username string, password string) (*model.User, error) {
	// Проверка инициализирован ли интерфейс UserService
	if r.UserService == nil {
		return nil, errors.New("user service is not initialized")
	}
	// Проверка существует ли пользователь уже в базе данных
	_, err := r.UserService.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	password = r.UserService.HashPassword(password)
	// Вызов инфраструктурной функции создания пользователя, которая прикреплена к интерфейсу сервиса пользователей
	createdUser, err := r.UserService.UserCreate(ctx, username, password)
	if err != nil {
		return nil, err
	}
	// Проверка успешности создания пользователя
	if createdUser == nil {
		return nil, errors.New("failed to create user")
	}
	user := model.User{
		ID:       createdUser.ID,
		Username: username,
		Password: password,
	}
	return &user, nil
}

// Функция редактирования имени полльзователя (защищено токеном)
func (r *mutationResolver) UpdateUserUsername(ctx context.Context, id string, username string) (*model.User, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET username=$2 WHERE id=$1", id, username)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Username: username}, nil
}

// Функция редактирования пароля (защищено токеном)
func (r *mutationResolver) UpdateUserPassword(ctx context.Context, id string, password string) (*model.User, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET password=$2 WHERE id=$1", id, password)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Password: password}, nil
}

// Функция удаления собственной учетной записи (защищено токеном)
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	authUser := middleware.CtxValue(ctx)
	if authUser == nil {
		return nil, errors.New("unauthorized")
	}
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, text, author_id FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Text, &post.AuthorID); err != nil {
			return nil, err
		}

		post.Author, err = r.UserService.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			return nil, err
		}

		post.Comments, err = r.UserService.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := r.DB.QueryRowContext(ctx, "SELECT id, text, author_id FROM post WHERE id=$1", id).Scan(&post.ID, &post.Text, &post.AuthorID)
	if err != nil {
		return nil, err
	}

	post.Author, err = r.UserService.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		return nil, err
	}

	post.Comments, err = r.UserService.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, text string, permissionToComment bool) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("create post not auth")
	}
	fmt.Println(user.ID)
	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO post (id, text, author_id, commentable) VALUES ($1, $2, $3, $4)", id, text, user.ID, permissionToComment)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: id, Text: text, AuthorID: user.ID}, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, text string) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := r.DB.ExecContext(ctx, "UPDATE post SET text=$2 WHERE id=$1", id, text)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: id, Text: text}, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	var post model.Post
	err := r.DB.QueryRowContext(ctx, "SELECT id, text, author_id FROM post WHERE id=$1", id).Scan(&post.ID, &post.Text, &post.AuthorID)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.ExecContext(ctx, "DELETE FROM post WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context) ([]*model.CommentResponse, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment")
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

		comment.Author, err = r.UserService.GetUserByID(ctx, comment.AuthorID)
		if err != nil {
			return nil, err
		}

		comment.Replies, err = r.UserService.GetCommentsByParentID(ctx, comment.ID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}
	return comments, nil
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*model.CommentResponse, error) {
	var comment model.CommentResponse
	err := r.DB.QueryRowContext(ctx, "SELECT id, comment, author_id, post_id, parent_comment_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.PostID, &comment.ParentCommentID)
	if err != nil {
		return nil, err
	}

	comment.Author, err = r.UserService.GetUserByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, err
	}

	comment.Replies, err = r.UserService.GetCommentsByParentID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, comment string, itemId string) (*model.CommentResponse, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	var commentAble bool
	err := r.DB.QueryRowContext(ctx, "SELECT commentable FROM post WHERE id=$1", itemId).Scan(&commentAble)
	if err != nil {
		return nil, errors.New("Author of this post turned off comments")
	}
	var isReply bool
	var parentCommentID *string
	var postID string
	err = r.DB.QueryRowContext(ctx, "SELECT post_id FROM comment WHERE id=$1", itemId).Scan(&postID)
	if err == sql.ErrNoRows {
		postID = itemId
		isReply = false		
	} else if err != nil{
		return nil, err
	} else {
		parentCommentID = &itemId
		isReply = true
	}
	id := uuid.New().String()
	var query string
	if isReply {
		query = "INSERT INTO comment (id, comment, author_id, post_id, parent_comment_id) VALUES ($1, $2, $3, $4, $5)"
		_, err := r.DB.ExecContext(ctx, query, id, comment, user.ID, postID, itemId)
		if err != nil {
			return nil, err
		}
		return &model.CommentResponse{ID: id, Comment: comment, AuthorID: user.ID, PostID: postID, ParentCommentID: parentCommentID}, nil
	} else {
		query = "INSERT INTO comment (id, comment, author_id, post_id, parent_comment_id) VALUES ($1, $2, $3, $4, NULL)"
		_, err := r.DB.ExecContext(ctx, query, id, comment, user.ID, itemId)
		if err != nil {
			return nil, err
		}
		return &model.CommentResponse{ID: id, Comment: comment, AuthorID: user.ID, PostID: postID}, nil
	}
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := r.DB.ExecContext(ctx, "UPDATE comment SET comment=$2 WHERE id=$1", id, comment)
	if err != nil {
		return nil, err
	}
	return &model.Comment{ID: id, Comment: comment}, nil
}

