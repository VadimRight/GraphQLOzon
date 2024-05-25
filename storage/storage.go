package storage

import "github.com/VadimRight/GraphQLOzon/graph/model"

type Storage interface {
    GetUserByUsername(username string) (*model.User, error)
    CreateUser(username, password string) (*model.User, error)
    GetPostsByUserID(userID string) ([]*model.Post, error)
    GetCommentsByPostID(postID string) ([]*model.CommentResponse, error)
    GetCommentsByParentID(parentID string) ([]*model.CommentResponse, error)
    GetCommentsByUserID(userID string) ([]*model.CommentResponse, error)
    GetUserByID(userID string) (*model.User, error)
}
