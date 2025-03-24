package services_auth

import (
	"github.com/gin-gonic/gin"
	models_auth "pashmak.com/pashmak/models"
)

func (as *AuthService) CheckUserPassword(email string, password string) bool{
	var user models_auth.User
	result := as.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return false
	}
	return user.CheckPassword(password)
}

func LoginWithPassword(email string, password string){

}