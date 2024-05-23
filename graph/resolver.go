
// resolver/resolver.go
package graph

import (
	"context"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"database/sql"
	"github.com/google/uuid"
)

type Resolver struct{
	DB *sql.DB
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, username, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, username string, password string) (*model.User, error) {
	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", id, username, password)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Name: username, Password: password}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, username string, password string) (*model.User, error) {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET username=$2, password=$3 WHERE id=$1", id, username, password)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Name: username, Password: password}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, text string, authorId string) (*model.Post, error) {
	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO post (id, text, author_id) VALUES ($1, $2, $3)", id, text, authorId)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: id, Text: text, AuthorID: authorId}, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, text string) (*model.Post, error) {
	_, err := r.DB.ExecContext(ctx, "UPDATE post SET text=$2 WHERE id=$1", id, text)
	if err != nil {
		return nil, err
	}
	return &model.Post{ID: id, Text: text}, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*model.Post, error) {
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
	rows, err := r.DB.QueryContext(ctx, "SELECT id, comment, author_id FROM comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.AuthorID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment
	err := r.DB.QueryRowContext(ctx, "SELECT id, comment, author_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.Comment, &comment.AuthorID)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, comment string, authorId string) (*model.Comment, error) {
	id := uuid.New().String()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO comment (id, comment, author_id) VALUES ($1, $2, $3)", id, comment, authorId)
	if err != nil {
		return nil, err
	}
	return &model.Comment{ID: id, Comment: comment, AuthorID: authorId}, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	_, err := r.DB.ExecContext(ctx, "UPDATE comment SET comment=$2 WHERE id=$1", id, comment)
	if err != nil {
		return nil, err
	}
	return &model.Comment{ID: id, Comment: comment}, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment
	err := r.DB.QueryRowContext(ctx, "SELECT id, comment, author_id FROM comment WHERE id=$1", id).Scan(&comment.ID, &comment.Comment, &comment.AuthorID)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.ExecContext(ctx, "DELETE FROM comment WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

