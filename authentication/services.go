package authentication

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/initializers"
	"pashmak.com/pashmak/models"
)

type AuthService struct{
  DB          *gorm.DB
  RedisClient *redis.Client
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client) *AuthService{
  return &AuthService{
    DB : db,
    RedisClient: redisClient,
  }
}

func GenerateOTP()string{
  rand.NewSource(time.Now().UnixNano())
  otp := rand.Intn(10000)
  return fmt.Sprintf("%04d", otp)
}

func (as *AuthService)CheckExistance(email string) (bool, error){
  var user models.User
  fmt.Println(as.DB)
  // FIXME: initializers.DB should be replaced with as.DB but it causes error
  result := initializers.DB.First(&user, "email = ?", email)
  
  if result.Error != nil{
    if result.Error == gorm.ErrRecordNotFound{
      return false, nil
    }
    return false, result.Error
  }
  return user.ID != 0, nil
}

func (as *AuthService)ValidateUser(email string) (bool, error){
  exists, err := as.CheckExistance(email)
  if err != nil{
    return false, fmt.Errorf("failed to check user existence: %w", err)
  }

  if exists{
    userOTP := GenerateOTP()
    ctx := context.Background()
    // FIXME: initializers.RedisClient should be replaced with as.RedisClient but it causes error
    err = initializers.RedisClient.Set(ctx, email, userOTP, 5*time.Minute).Err()
    if err != nil{
      return true, fmt.Errorf("failed to set OTP in Redis: %w", err)
    }
    return true, nil
  }
  return false, nil
}