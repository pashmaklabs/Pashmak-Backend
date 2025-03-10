package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) StartEmailAuth(c *gin.Context) {
	// Read body
	var body StartEmailAuthRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "error reading request body",
			"errorCode": "INVALID_REQUEST_BODY",
		})
		return
	}

	// Pass to service
	resp := ac.authService.ValidateUser(body.Email)
	
	// Send response
	if !resp{
	  c.JSON(http.StatusNotFound, gin.H{
		"status": "error",
		"message": "User not found",
		"errorCode": "USER_NOT_FOUND",
	  })
	  return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "User found",
	})
}
