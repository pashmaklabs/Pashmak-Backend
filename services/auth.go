package services_auth

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"pashmak.com/pashmak/models"
)

type AuthService struct{
  DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService{
  return &AuthService{db}
}

func GnerateOTP()string{
  rand.NewSource(time.Now().UnixNano())
  otp := rand.Intn(10000)
  return fmt.Sprintf("%04d", otp)
}

func (as *AuthService)CheckExistance(email string) bool{
  var user models_auth.User
  as.DB.First(&user, "email = ?", email)
  return user.ID != 0
}

func (as *AuthService)ValidateUser(email string) bool{
  if as.CheckExistance(email){
    userotp := GnerateOTP()
    fmt.Println(userotp)
    // TODO: store OTP in redis
    return true
  }
  return false
}