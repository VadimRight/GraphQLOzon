// resolver/resolver.go
package graph

import (
	"context"
	"errors"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/google/uuid"
	"fmt"
)

// Тип Resolver, который ответственен за работу с данными в нашей схеме GraphQL
type Resolver struct{
	UserService service.UserService
	CommentService service.CommentService
	PostService service.PostService
}

// Функция возвращающая тип Запросов нашего резольвера
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Функция возвращающая тип Мутаций нашего резольвера
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Типы используемых методов GraphQL - тип запросов (аналог GET) и мутации (запросы, способных изменить данные)
type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Функция получения всех пользователей
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.UserService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Получения пользователя по его ID
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, err := r.UserService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Функция получения пользователя по его нику
func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := r.UserService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
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

// Метод получения всех постов
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.PostService.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		// Заполняем автора поста
		post.AuthorPost, err = r.UserService.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для поста
		post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID)
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
			comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID)
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
	return posts, nil
}

// Метод получения постов по ID пользователя
func (r *queryResolver) PostsByUserID(ctx context.Context, userID string) ([]*model.Post, error) {
	posts, err := r.PostService.GetPostsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		// Заполняем автора поста
		post.AuthorPost, err = r.UserService.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для поста
		post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID)
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
			comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID)
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
	return posts, nil
}

// Метод получения поста по его ID
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	post, err := r.PostService.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Заполняем автора поста
	post.AuthorPost, err = r.UserService.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		return nil, err
	}

	// Получаем комментарии для поста
	post.Comments, err = r.CommentService.GetCommentsByPostID(ctx, post.ID)
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
		comment.Replies, err = r.CommentService.GetCommentsByParentID(ctx, comment.ID)
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

	return post, nil
}

// Метод создания поста
func (r *mutationResolver) CreatePost(ctx context.Context, text string, permissionToComment bool) (*model.Post, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	fmt.Println("Creating comment for user ID:", user.ID)
	id := uuid.New().String()
	post, err := r.PostService.CreatePost(ctx, id, text, user.ID, permissionToComment)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Метод получения всех коментариев
func (r *queryResolver) Comments(ctx context.Context, limit *int, offset *int) ([]*model.CommentResponse, error) {
	comments, err := r.CommentService.GetAllComments(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Метод получения комментария по его ID
func (r *queryResolver) Comment(ctx context.Context, id string) (*model.CommentResponse, error) {
	comment, err := r.CommentService.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// Метод создания комментария
func (r *mutationResolver) CreateComment(ctx context.Context, commentText string, itemId string) (*model.CommentResponse, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	fmt.Println("Creating comment for user ID:", user.ID)
	comment, err := r.CommentService.CreateComment(ctx, commentText, itemId, user.ID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
