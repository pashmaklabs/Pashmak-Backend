package services_auth

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
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

func SendMail(Email string, userOTP string) error {
	from := bootstrap.EMAIL_ADDR
	password := bootstrap.EMAIL_PASSWORD
	smtpHost := bootstrap.EMAIL_HOST
	smtpPort := bootstrap.EMAIL_PORT

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
		return err
	}
	return nil
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

func (as *AuthService) StoreOTPAndSendEmail(email string) error {
	userOTP := GenerateOTP()
	ctx := context.Background()
	if err := as.RedisClient.Set(ctx, email, userOTP, 2*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}
	if err := SendMail(email, userOTP); err != nil {
		return err
	}
	return nil
}

func (as *AuthService) ValidateUser(email string) (bool, error) {
	exists, err := as.CheckExistance(email)
	if err != nil {
		return exists, fmt.Errorf("failed to check user existence: %w", err)
	}

	if err := as.StoreOTPAndSendEmail(email); err != nil {
		return exists, err
	}

	return exists, nil
}

func (as *AuthService) ValidateOTP(Email string, RecievedOTP string) (bool, error) {
	ctx := context.Background()
	realOTP, err := as.RedisClient.Get(ctx, Email).Result()
	if err != nil {
		return false, fmt.Errorf("failed to get OTP from redis: %w", err)
	}

	if realOTP != RecievedOTP {
		return false, nil
	}
	return true, nil
}
