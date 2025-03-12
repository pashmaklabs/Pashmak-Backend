package authentication

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/initializers"
	"pashmak.com/pashmak/models"
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

func SendMail(Email string, userOTP string) {
	from := "pashmak471@gmail.com"
	password := initializers.EMAIL_PASSWORD

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	htmlContent := fmt.Sprintf(`
	<html>
		<body>
			<h1>Your Verification Code:</h2>
			<h1 style="font-size: 36px; color: #007BFF;">%s</h1>
			<p>Please use this code to verify your email address.</p>
		</body>
	</html>
	`, userOTP)

	message := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"+
			"%s",
		strings.Join([]string{Email}, ","),
		"Verify Email",
		htmlContent,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{Email}, message)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GenerateOTP() string {
	rand.NewSource(time.Now().UnixNano())
	otp := rand.Intn(10000)
	return fmt.Sprintf("%04d", otp)
}

func (as *AuthService) CheckExistance(Email string) (bool, error) {
	var user models.User
	// fmt.Println(as.DB)
	// FIXME: initializers.DB should be replaced with as.DB but it causes error
	result := initializers.DB.First(&user, "email = ?", Email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return user.ID != 0, nil
}

func (as *AuthService) ValidateUser(Email string) (bool, error) {
	exists, err := as.CheckExistance(Email)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		userOTP := GenerateOTP()
		ctx := context.Background()
		// FIXME: initializers.RedisClient should be replaced with as.RedisClient but it causes error
		err = initializers.RedisClient.Set(ctx, Email, userOTP, 5*time.Minute).Err()
		if err != nil {
			return true, fmt.Errorf("failed to set OTP in Redis: %w", err)
		}
		SendMail(Email, userOTP)
		return true, nil
	}
	return false, nil
}

func (as *AuthService)ValidateOTP(Email string, RecievedOTP string)(bool, error){
  ctx := context.Background()
  realOTP, err := initializers.RedisClient.Get(ctx, Email).Result()
  if err != nil {
    return false, fmt.Errorf("failed to get OTP from redis: %w", err)
  }

  if realOTP != RecievedOTP{
    return false, nil
  }
  return true, nil
}
