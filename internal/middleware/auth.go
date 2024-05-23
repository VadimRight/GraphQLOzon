// middlewares/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/VadimRight/GraphQLOzon/internal/service"
)

type authString string

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.Next()
			return
		}

		bearer := "Bearer "
		if !strings.HasPrefix(auth, bearer) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		token := auth[len(bearer):]

		validate, err := service.JwtValidate(context.Background(), token)
		if err != nil || !validate.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		customClaim, _ := validate.Claims.(*service.JwtCustomClaim)
		ctx := context.WithValue(c.Request.Context(), authString("auth"), customClaim)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	return raw
}
