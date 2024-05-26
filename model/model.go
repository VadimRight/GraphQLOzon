package model

// Comment представляет собой структуру комментария
type Comment struct {
	ID            string `json:"id"`
	Comment       string `json:"comment"`
	AuthorID      string `json:"authorId"`
	ItemID        string `json:"itemId"`
	AuthorComment *User  `json:"authorComment"`
}

// CommentResponse представляет собой структуру ответа на комментарий
type CommentResponse struct {
	ID              string             `json:"id"`
	Comment         string             `json:"comment"`
	AuthorID        string             `json:"authorId"`
	PostID          string             `json:"postId"`
	ParentCommentID *string            `json:"parentCommentID,omitempty"`
	AuthorComment   *User              `json:"authorComment"`
	Replies         []*CommentResponse `json:"replies"`
}

// Mutation представляет собой структуру мутаций GraphQL
type Mutation struct {
}

// Post представляет собой структуру поста
type Post struct {
	ID          string             `json:"id"`
	Text        string             `json:"text"`
	AuthorID    string             `json:"authorId"`
	AuthorPost  *User              `json:"authorPost"`
	Comments    []*CommentResponse `json:"comments"`
	Commentable bool               `json:"commentable"`
}

// Query представляет собой структуру запросов GraphQL
type Query struct {
}

// Token представляет собой структуру токена
type Token struct {
	Token string `json:"token"`
}

// User представляет собой структуру пользователя
type User struct {
	ID       string             `json:"id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Posts    []*Post            `json:"posts"`
	Comments []*CommentResponse `json:"comments"`
}
