package graph

import (
	"context"
	"errors"

	"github.com/VadimRight/GraphQLOzon/internal/middleware"
	"github.com/VadimRight/GraphQLOzon/model"
	"github.com/google/uuid"
)

// Метод получения всех постов
func (r *queryResolver) Posts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	posts, err := r.PostUsecase.GetAllPosts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		// Заполняем автора поста
		post.AuthorPost, err = r.UserUsecase.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для поста
		post.Comments, err = r.CommentUsecase.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, comment := range post.Comments {
			// Заполняем автора комментария
			comment.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, comment.AuthorID)
			if err != nil {
				return nil, err
			}

			// Получаем ответы для каждого комментария
			comment.Replies, err = r.CommentUsecase.GetCommentsByParentID(ctx, comment.ID, limit, offset)
			if err != nil {
				return nil, err
			}

			for _, reply := range comment.Replies {
				// Заполняем автора ответа
				reply.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, reply.AuthorID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return posts, nil
}

// Метод получения постов по ID пользователя
func (r *queryResolver) PostsByUserID(ctx context.Context, userID string, limit, offset *int) ([]*model.Post, error) {
	posts, err := r.PostUsecase.GetPostsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		// Заполняем автора поста
		post.AuthorPost, err = r.UserUsecase.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для поста
		post.Comments, err = r.CommentUsecase.GetCommentsByPostID(ctx, post.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, comment := range post.Comments {
			// Заполняем автора комментария
			comment.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, comment.AuthorID)
			if err != nil {
				return nil, err
			}

			// Получаем ответы для каждого комментария
			comment.Replies, err = r.CommentUsecase.GetCommentsByParentID(ctx, comment.ID, limit, offset)
			if err != nil {
				return nil, err
			}

			for _, reply := range comment.Replies {
				// Заполняем автора ответа
				reply.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, reply.AuthorID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return posts, nil
}

// Метод получения поста по его ID
func (r *queryResolver) Post(ctx context.Context, id string, limit, offset *int) (*model.Post, error) {
	post, err := r.PostUsecase.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Заполняем автора поста
	post.AuthorPost, err = r.UserUsecase.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		return nil, err
	}

	// Получаем комментарии для поста
	post.Comments, err = r.CommentUsecase.GetCommentsByPostID(ctx, post.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, comment := range post.Comments {
		// Заполняем автора комментария
		comment.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, comment.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получаем ответы для каждого комментария
		comment.Replies, err = r.CommentUsecase.GetCommentsByParentID(ctx, comment.ID, limit, offset)
		if err != nil {
			return nil, err
		}

		for _, reply := range comment.Replies {
			// Заполняем автора ответа
			reply.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, reply.AuthorID)
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
	id := uuid.New().String()
	post, err := r.PostUsecase.CreatePost(ctx, id, text, user.ID, permissionToComment)
	if err != nil {
		return nil, err
	}
	return post, nil
}
