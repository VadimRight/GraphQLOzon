package graph

import (
	"context"
	"errors"
	"github.com/VadimRight/GraphQLOzon/graph/model"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
)

// Метод получения всех комментариев
func (r *queryResolver) Comments(ctx context.Context, limit, offset *int) ([]*model.CommentResponse, error) {
	comments, err := r.CommentUsecase.GetAllComments(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
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

	return comments, nil
}

// Метод получения комментария по его ID
func (r *queryResolver) Comment(ctx context.Context, id string, limit, offset *int) (*model.CommentResponse, error) {
	comment, err := r.CommentUsecase.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Заполняем автора комментария
	comment.AuthorComment, err = r.UserUsecase.GetUserByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, err
	}

	// Получаем ответы для комментария с поддержкой пагинации
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

	return comment, nil
}

// Метод создания комментария
func (r *mutationResolver) CreateComment(ctx context.Context, commentText string, itemId string) (*model.CommentResponse, error) {
	user := middleware.CtxValue(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}
	comment, err := r.CommentUsecase.CreateComment(ctx, commentText, itemId, user.ID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
