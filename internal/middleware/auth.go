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

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// Middleware для аутентификации
func (a *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", auth)

		// Если заголовок пустой, пропускаем запрос дальше
		if auth == "" {
			c.Next()
			return
		}

		bearer := "Bearer "
		// Проверяем, начинается ли заголовок с "Bearer "
		if !strings.HasPrefix(auth, bearer) {
			// Если нет, возвращаем ошибку 403 Forbidden
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		// Извлекаем сам токен из заголовка
		token := auth[len(bearer):]
		fmt.Println("Token:", token)

		// Валидируем токен с помощью службы JWT
		validate, err := a.authService.ValidateToken(context.Background(), token)
		if err != nil || !validate.Valid {
			fmt.Println("Validation Error:", err)
			// Если токен недействителен, возвращаем ошибку 403 Forbidden
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		// Извлекаем кастомные claims из токена
		customClaim, _ := validate.Claims.(*service.JwtCustomClaim)
		// Добавляем claims в контекст запроса
		ctx := context.WithValue(c.Request.Context(), AuthKey, customClaim)
		c.Request = c.Request.WithContext(ctx)

		// Пропускаем запрос дальше
		c.Next()
	}
}

// Функция для получения значений из контекста
func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	// Извлекаем кастомные claims из контекста
	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	return raw
}
