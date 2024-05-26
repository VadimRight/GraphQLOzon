package graph

import (
	"context"
	"errors"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/service"
)

// Функция получения всех пользователей
func (r *queryResolver) Users(ctx context.Context, limit, offset *int) ([]*model.User, error) {
	users, err := r.UserService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		// Получаем посты для пользователя
		user.Posts, err = r.PostService.GetPostsByUserID(ctx, user.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, post := range user.Posts {
			// Заполняем автора поста
			post.AuthorPost, err = r.UserService.GetUserByID(ctx, post.AuthorID)
			if err != nil {
				return nil, err
			}

			// Получаем комментарии для поста
			post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID, limit, offset)
			if err != nil {
				return nil, err
			}

			for _, comment := range post.Comments {
				// Заполняем автора комментария
				comment.AuthorComment, err = r.UserService.GetUserByID(ctx, comment.AuthorID)
				if err != nil {
					return nil, err
				}

				// Получаем ответы для каждого комментария
				comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID, limit, offset)
				if err != nil {
					return nil, err
				}

				for _, reply := range comment.Replies {
					// Заполняем автора ответа
					reply.AuthorComment, err = r.UserService.GetUserByID(ctx, reply.AuthorID)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return users, nil
}

// Получения пользователя по его ID
func (r *queryResolver) User(ctx context.Context, id string, limit, offset *int) (*model.User, error) {
	user, err := r.UserService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Получаем посты для пользователя
	user.Posts, err = r.PostService.GetPostsByUserID(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, post := range user.Posts {
		// Заполняем автора поста
		post.AuthorPost, err = r.UserService.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для поста
		post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, comment := range post.Comments {
			// Заполняем автора комментария
			comment.AuthorComment, err = r.UserService.GetUserByID(ctx, comment.AuthorID)
			if err != nil {
				return nil, err
			}

			// Получаем ответы для каждого комментария
			comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID, limit, offset)
			if err != nil {
				return nil, err
			}

			for _, reply := range comment.Replies {
				// Заполняем автора ответа
				reply.AuthorComment, err = r.UserService.GetUserByID(ctx, reply.AuthorID)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return user, nil
}

func (r *queryResolver) UserByUsername(ctx context.Context, username string, limit, offset *int) (*model.User, error) {
	user, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Получение постов пользователя с учетом пагинации
	user.Posts, err = r.PostService.GetPostsByUserID(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, post := range user.Posts {
		// Заполнение автора поста
		post.AuthorPost, err = r.UserService.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получение комментариев для поста с учетом пагинации
		post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, comment := range post.Comments {
			// Заполнение автора комментария
			comment.AuthorComment, err = r.UserService.GetUserByID(ctx, comment.AuthorID)
			if err != nil {
				return nil, err
			}

			// Получение ответов для каждого комментария с учетом пагинации
			comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID, limit, offset)
			if err != nil {
				return nil, err
			}

			for _, reply := range comment.Replies {
				// Заполнение автора ответа
				reply.AuthorComment, err = r.UserService.GetUserByID(ctx, reply.AuthorID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return user, nil
}

// Метод логина пользователя
func (r *mutationResolver) LoginUser(ctx context.Context, username string, password string) (*model.Token, error) {
	getUser, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
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

// Метод регистрации пользователя
func (r *mutationResolver) RegisterUser(ctx context.Context, username string, password string) (*model.User, error) {
	_, err := r.UserService.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	password = r.UserService.HashPassword(password)
	createdUser, err := r.UserService.UserCreate(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

