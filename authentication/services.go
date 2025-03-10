package authentication

import (
  "pashmak.com/pashmak/initializers"
  "pashmak.com/pashmak/models"
  "gorm.io/gorm"
  "math/rand"
  "time"
  "fmt"
)

type AuthService struct{
  DB *gorm.DB
}

func NewAuthService() *AuthService{
  return &AuthService{}
}

func GnerateOTP()string{
  rand.NewSource(time.Now().UnixNano())
  otp := rand.Intn(10000)
  return fmt.Sprintf("%04d", otp)
}

func CheckExistance(email string) bool{
  var user models.User
  initializers.DB.First(&user, "email = ?", email)
  return user.ID != 0
}

func (as *AuthService)ValidateUser(email string) bool{
  if CheckExistance(email){
    userotp := GnerateOTP()
    fmt.Println(userotp)
    // TODO: store OTP in redis
    return true
  }
  return false
}