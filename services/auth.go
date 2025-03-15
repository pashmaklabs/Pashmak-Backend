package services_auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	models_auth "pashmak.com/pashmak/models"
)

type AuthService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client) *AuthService {
	return &AuthService{
		DB:          db,
		RedisClient: redisClient,
	}
}

func GenerateOTP() string {
	rand.NewSource(time.Now().UnixNano())
	otp := rand.Intn(10000)
	return fmt.Sprintf("%04d", otp)
}

func (as *AuthService) CheckExistance(email string) (bool, error) {
	var user models_auth.User
	result := as.DB.First(&user, "email = ?", email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (as *AuthService)StoreInRedis(email string) error{
  userOTP := GenerateOTP()
  ctx := context.Background()
  if err := as.RedisClient.Set(ctx, email, userOTP, 2*time.Minute).Err(); err != nil{
    return fmt.Errorf("failed to store OTP in Redis: %w", err)
  }
  return nil
}

func (as *AuthService) ValidateUser(email string) (bool, error) {
	exists, err := as.CheckExistance(email)
	if err != nil {
		return exists, fmt.Errorf("failed to check user existence: %w", err)
	}

  if err := as.StoreInRedis(email); err != nil {
		return exists, err
	}

  return exists, nil
}
