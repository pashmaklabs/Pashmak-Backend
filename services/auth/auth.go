package services_auth

import (
	"github.com/gin-gonic/gin"
	models_auth "pashmak.com/pashmak/models"
)

func (as *AuthService) GetUserByGmail(email string) (models_auth.User, error){
	var user models_auth.User
	result := as.DB.First(&user, "email = ?", email)
	return user, result.Error
}

func (as *AuthService) CreateUser(email string) error{
	result := as.DB.Create(&models_auth.User{Email: email})
	return result.Error
}

func (as *AuthService) IsUserLoggedIn(c *gin.Context) bool{
	_, exists := c.Get("user")
	return exists
}