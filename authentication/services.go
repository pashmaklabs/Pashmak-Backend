package authentication

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	// "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/initializers"
	"pashmak.com/pashmak/models"
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

func (as *AuthService)ValidateUser(email string) (bool, error){
  if CheckExistance(email){
    userotp := GnerateOTP()
    fmt.Println(userotp)
    ctx := context.Background()
    err := initializers.RedisClient.Set(ctx, email, userotp, 5*time.Minute).Err()
    if err != nil{
      return false, err
    }
    // TODO: store OTP in redis
    return true, nil
  }
  return false, nil
}