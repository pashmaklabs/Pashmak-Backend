package middlewares_auth

import (
	"net/http"

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
		token, err := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"message": "ابتدا باید وارد شوید",
			})
			c.Abort()
			return
		} else {
			claim, err := am.authService.VerifyJWT(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": "error",
					"message": "در ورود مشکلی پیش آمده",
				})
				c.Abort()
				return
			}
			c.Set("user", claim.UserInfo) // Needs consideration
			c.Set("claim", &claim)
			c.Next()
		}
	}
}