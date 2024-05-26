package api

import (
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/VadimRight/GraphQLOzon/graph"
	"github.com/VadimRight/GraphQLOzon/storage"
	"github.com/VadimRight/GraphQLOzon/internal/config"
	"log"
	"github.com/VadimRight/GraphQLOzon/internal/service"
	"github.com/VadimRight/GraphQLOzon/internal/middleware"
	"github.com/VadimRight/GraphQLOzon/internal/usecase"
)

// Функция инициализации сервера
func InitServer(cfg *config.Config, storage storage.Storage) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()

	// Инициализация сервисов и middleware
	authService := service.NewAuthService()
	authMiddleware := middleware.NewAuthMiddleware(authService)

	r.Use(authMiddleware.Handler())

	r.POST("/graphql", graphqlHandler(storage, authService))
	r.GET("/", playgroundHandler())

	log.Println("connect to http://localhost:8000/ for GraphQL playground")
	log.Fatal(r.Run(":8000"))
}

// Хэндлер для непосредственно нашей схемы GraphQL
func graphqlHandler(storage storage.Storage, authService service.AuthService) gin.HandlerFunc {
	postUsecase := usecase.NewPostUsecase(storage)
	commentUsecase := usecase.NewCommentUsecase(storage)
	userUsecase := usecase.NewUserUsecase(storage, commentUsecase, service.NewPasswordService(), authService)

	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			UserUsecase:    userUsecase,
			PostUsecase:    postUsecase,
			CommentUsecase: commentUsecase,
		},
	}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Хендлер для песочницы, где можно отправлять HTTP запросы от клиента на сервер
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
