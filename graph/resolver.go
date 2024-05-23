
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
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Username, &user.Password)
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
	return &model.User{ID: id, Username: username, Password: password}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, username string, password string) (*model.User, error) {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET username=$2, password=$3 WHERE id=$1", id, username, password)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Username: username, Password: password}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
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
