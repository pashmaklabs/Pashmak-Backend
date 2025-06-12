package middlewares_auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models/auth"
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
func (am *AuthMiddleware) LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("pashmak_authentication")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "ابتدا باید وارد شوید",
			})
			return
		} else {
			claim, err := am.authService.VerifyJWT(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
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

func HasPermission(db *gorm.DB, userID uint, permissionName string) bool {
	var user models_auth.User
	if err := db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return false
	}

	for _, perm := range user.Role.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}

func (am *AuthMiddleware)PermissionMiddleware(db *gorm.DB, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}
		userinfo := user.(services_auth.UserInfo)

		if !HasPermission(db, userinfo.ID, permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()

		// Add middleware in this format: PermissionMiddleware(db, "create_post")
	}
}

func (am *AuthMiddleware) AuthOrAnonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("pashmak_authentication")
		if err == nil && token != "" {
			claim, err := am.authService.VerifyJWT(token)
			if err == nil {
				c.Set("user", *claim.UserInfo)
				c.Set("claim", &claim)
			} else {
				// Invalid JWT, treating as anonymous
			}
		}
		c.Next()
	}
}
