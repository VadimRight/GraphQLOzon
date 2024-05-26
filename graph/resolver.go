package graph

import (
	"github.com/VadimRight/GraphQLOzon/internal/service"
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

