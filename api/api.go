package api


import (
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/VadimRight/GraphQLOzon/graph"
	"github.com/VadimRight/GraphQLOzon/bootstrap"
	"log"
)

func InitServer(cfg *bootstrap.Config, storage *bootstrap.Storage) {	
	r := gin.Default()
	r.POST("/query", graphqlHandler(storage))
	r.GET("/", playgroundHandler())
	r.Run(cfg.Server.ServerAddress)
	log.Println("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(r.Run(":8080"))
}

func graphqlHandler(storage *bootstrap.Storage) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: storage.DB}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
