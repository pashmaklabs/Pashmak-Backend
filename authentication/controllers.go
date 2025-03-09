package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}


func (ac *AuthController)StartEmailAuth(c *gin.Context){
	var body StartEmailAuthRequest

	if c.Bind(body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "error reading request body",
		})
	}

	// Pass to service
	//resp := ac.authService.ValidateUser(body.Email)


	// Send response
}