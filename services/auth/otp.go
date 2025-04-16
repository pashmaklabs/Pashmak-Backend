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
	models_auth "pashmak.com/pashmak/models/auth"
)


// TODO: move basic services_auth setup to another file(not in otp!)
type AuthService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	AppConfig   *bootstrap.AppConfig
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client, appConfig *bootstrap.AppConfig) *AuthService {
	return &AuthService{
		DB:          db,
		RedisClient: redisClient,
		AppConfig:   appConfig,
	}
}

func (as *AuthService) GnerateOTP()string{
  rand.NewSource(time.Now().UnixNano())
  otp := rand.Intn(10000)
  return fmt.Sprintf("%04d", otp)
}

func (as *AuthService) SendMail(Email string, userOTP string) error {
	from := as.AppConfig.EmailAddr
	password := as.AppConfig.EmailPassword
	smtpHost := as.AppConfig.EmailHost
	smtpPort := as.AppConfig.EmailPort

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

func (as *AuthService) CheckExistance(email string) error {
	var user models_auth.User
	result := as.DB.First(&user, "email = ?", email)
	return result.Error
}

func (as *AuthService) StoreOTPAndSendEmail(email string) error {
	userOTP := GenerateOTP()
	ctx := context.Background()
	if err := as.RedisClient.Set(ctx, email, userOTP, 2*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}
	if err := as.SendMail(email, userOTP); err != nil {
		return err
	}
	return nil
}

func (as *AuthService) ValidateUser(email string) (bool, error) {
	err := as.CheckExistance(email)
	if err != nil &&  err != gorm.ErrRecordNotFound{
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	if err := as.StoreOTPAndSendEmail(email); err != nil {
		return false, err
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

func (as *AuthService) ValidateOTP(Email string, RecievedOTP string) (bool, error) {
	ctx := context.Background()
	realOTP, err := as.RedisClient.Get(ctx, Email).Result()
	if err != nil {
		if err == redis.Nil{
			return false, err
		}
		return false, fmt.Errorf("failed to get OTP from redis: %w", err)
	}

	if realOTP != RecievedOTP {
		return false, nil
	}
	return true, nil
}
