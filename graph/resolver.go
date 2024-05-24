// resolver/resolver.go
package graph

import (
	"context"
	"database/sql"
	"errors"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
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

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	rows, err := r.DB.QueryContext(ctx, "SELECT id, username, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	authUser := middleware.CtxValue(ctx)
	if authUser == nil {
		return nil, errors.New("unauthorized")
	}

	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, username string, password string) (*model.Token, error) {
	getUser, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			return nil, postgresErr
		}
		return nil, err
	}
	if getUser.Password != password {
		return nil, errors.New("Password incorrect")
	}
	token, err := service.JwtGenerate(ctx, getUser.ID)
	if err != nil {
		return nil, err
	}
	return &model.Token{Token: token}, nil
}

func (r *mutationResolver) RegisterUser(ctx context.Context, username string, password string) (*model.User, error) {
	_, err := r.UserService.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	createdUser, err := r.UserService.UserCreate(ctx, username, password)
	return &model.User{ID: createdUser.ID, Username: username, Password: password}, nil
}

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
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

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
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	var post model.Post
	err := r.DB.QueryRowContext(ctx, "SELECT id, text, author_id FROM post WHERE id=$1", id).Scan(&post.ID, &post.Text, &post.AuthorID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, text string, authorId string) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO post (id, text, author_id) VALUES ($1, $2, $3)", id, text, authorId)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: id, Text: text, AuthorID: authorId}, nil
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

func (r *queryResolver) Comments(ctx context.Context) ([]*model.Comment, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	rows, err := r.DB.QueryContext(ctx, "SELECT id, comment, author_id, item_id FROM comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.ItemID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*model.Comment, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	var comment model.Comment
	err := r.DB.QueryRowContext(ctx, "SELECT id, comment, author_id, item_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.ItemID)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, comment string, authorId string, itemId string) (*model.Comment, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO comment (id, comment, author_id, item_id) VALUES ($1, $2, $3, $4)", id, comment, authorId, itemId)
	if err != nil {
		return nil, err
	}
	return &model.Comment{ID: id, Comment: comment, AuthorID: authorId, ItemID: itemId}, nil
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

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*model.Comment, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	var comment model.Comment
	err := r.DB.QueryRowContext(ctx, "SELECT id, comment, author_id, item_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.Comment, &comment.AuthorID, &comment.ItemID)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.ExecContext(ctx, "DELETE FROM comment WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
