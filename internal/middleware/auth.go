// middlewares/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/VadimRight/GraphQLOzon/internal/service"
)

type authString string

var AuthKey = authString("auth")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", auth)

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
		fmt.Println("Token:", token)

		validate, err := service.JwtValidate(context.Background(), token)
		if err != nil || !validate.Valid {
			fmt.Println("Validation Error:", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		customClaim, _ := validate.Claims.(*service.JwtCustomClaim)
		ctx := context.WithValue(c.Request.Context(), AuthKey, customClaim)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	return raw
}

