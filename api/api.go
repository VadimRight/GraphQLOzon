package api

import (
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/VadimRight/GraphQLOzon/graph"
	"github.com/VadimRight/GraphQLOzon/storage"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"log"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
)

// Функция инициализации сервера
func InitServer(cfg *bootstrap.Config, storage *storage.PostgresStorage) {	
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.POST("/graphql", graphqlHandler(storage))
	r.GET("/", playgroundHandler())

	log.Println("connect to http://localhost:8000/ for GraphQL playground")
	log.Fatal(r.Run(":8000"))
}

// Хэндлер для непосредственно нашей схемы GraphQL
func graphqlHandler(storage *storage.PostgresStorage) gin.HandlerFunc {
	userService := service.NewUserService(storage)
	postService := service.NewPostService(storage)
	commentService := service.NewCommentService(storage)
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{UserService: userService, PostService: postService, CommentService: commentService}}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Хендлер для песочницы, где можно отправлять HTTP запроса от клиента на сервер
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
