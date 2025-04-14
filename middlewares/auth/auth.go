package middlewares_auth

import (
	"log"
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
// TODO: When someone is logged in, he can't login again using auth endpoints
// TODO: When someone is not logged in, he can't logout
func (am *AuthMiddleware)LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("pashmak_authentication")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"message": "ابتدا باید وارد شوید",
			})
			return
		} else {
			claim, err := am.authService.VerifyJWT(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status": "error",
					"message": "در ورود مشکلی پیش آمده",
				})
				log.Println(err.Error())
				return
			}
			c.Set("user", *claim.UserInfo) // Needs consideration
			// [FIXME] : Why pass reference?
			c.Set("claim", &claim)
			c.Next()
		}
	}
}