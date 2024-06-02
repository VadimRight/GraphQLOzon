package graph

import (
	"github.com/VadimRight/GraphQLOzon/internal/usecase"
)

// Тип Resolver, который ответственен за работу с данными в нашей схеме GraphQL
type Resolver struct {
	UserUsecase    usecase.UserUsecase
	CommentUsecase usecase.CommentUsecase
	PostUsecase    usecase.PostUsecase
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
