package services_auth

import (
	"fmt"
	"math/rand"
	"time"
	"gorm.io/gorm"
)

type AuthService struct{
  DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService{
  return &AuthService{db}
}

func (as *AuthService) GnerateOTP()string{
  rand.NewSource(time.Now().UnixNano())
  otp := rand.Intn(10000)
  return fmt.Sprintf("%04d", otp)
}

func (as *AuthService) ValidateUser(email string) bool{
  if _, err := as.GetUserByGmail(email); err == nil{
    userotp := as.GnerateOTP()
    fmt.Println(userotp)
    // TODO: store OTP in redis
    return true
  }
  return false
}