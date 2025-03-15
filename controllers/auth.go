package controlllers_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serializers_auth "pashmak.com/pashmak/serializers"
	services_auth "pashmak.com/pashmak/services"
)

type AuthController struct {
	authService *services_auth.AuthService
}

func NewAuthController(authService *services_auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	// Read body
	var body serializers_auth.StartEmailAuthRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "error reading request body",
			"errorCode": "INVALID_REQUEST_BODY",
		})
		return
	}

	// Pass to service
	resp := ac.authService.ValidateUser(body.Email)

	// Send response
	if !resp {
		c.JSON(http.StatusNotFound, gin.H{
			"status":    "error",
			"message":   "کاربر یافت نشد",
			"errorCode": "USER_NOT_FOUND",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "کاربر یافت شد",
		"errorCode": "USER_FOUND",
	})
}
