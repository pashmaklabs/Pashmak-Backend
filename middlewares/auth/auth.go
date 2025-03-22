package middlewares_auth

import (
	"github.com/gin-gonic/gin"
	services_auth "pashmak.com/pashmak/services/auth"
)

type AuthMiddleware struct {
	authService *services_auth.AuthService
}

func NewAuthMiddleware(authService *services_auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (am *AuthMiddleware)LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "هدر authorization نیاز است."})
			c.Abort()
			return
		} else {
			claim, err := am.authService.VerifyJWT(token)
			if err != nil {
				c.JSON(401, gin.H{"error": "توکن نامعتبر است."})
				c.Abort()
				return
			}
			c.Set("user", &claim.UserInfo)
			c.Set("claim", &claim)
			c.Next()
		}
	}
}