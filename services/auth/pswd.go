package services_auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/redis/go-redis/v9"
	models_auth "pashmak.com/pashmak/models"
)

func (as *AuthService)SendResetPasswordMail(email string, userOTP string) error {
	err := as.SendMail(email, userOTP)
	return err
}

func (as *AuthService) SetUserPassword(user *models_auth.User, newpassword string) error {
	// TODO: Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(newpassword), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	
	// user.Password = newpassword
	result := as.DB.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (as *AuthService) CheckUserPassword(email string, newpassword string) (*models_auth.User, error) {
	var user models_auth.User
	result := as.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.Password == "" {
		return nil, errors.New("user has no password")
	}
	// TODO: Hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newpassword)); err != nil {
		return nil, err
	}
	if user.Password != newpassword {
		return nil, nil
	}
	return &user, nil
}

func (as *AuthService) LoginWithPassword(email string, password string) (string, error) {
	user, err := as.GetUserByGmail(email)
	if err != nil {
		return "", err
	}
	if user.Password == "" {
		return "", errors.New("user has no password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	jwt, err := as.GenerateJWT(user)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (as *AuthService) ForgetPassword(email string) error {
	userOTP := GenerateOTP()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	if err := as.RedisClient.Set(ctx, email, userOTP, 2*time.Minute).Err(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
            return errors.New("operation timed out")
        }
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}
	if err := as.SendResetPasswordMail(email, userOTP); err != nil {
		return err
	}

	return nil
}

func (as *AuthService) VerifyForgetPassword(email string, otp string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // Ensures resources are cleaned up
	realOTP, err := as.RedisClient.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil { // This format can be used instead of errors.Is(err, redis.Nil)
			return "", false, errors.New("OTP expired")
		}
		if ctx.Err() == context.DeadlineExceeded {
            return "", false, errors.New("operation timed out")
        }
		return "", false, fmt.Errorf("failed to get OTP from redis: %w", err)
	}

	if realOTP != otp {
		return "", false, nil
	}
	user, err := as.GetUserByGmail(email)
	if err != nil {
		return "", false, fmt.Errorf("failed to get user by email: %w", err)
	}
	jwt, err := as.GenerateJWT(user)
	if err != nil {
		return "", false, fmt.Errorf("failed to generate JWT: %w", err)
	}
	return jwt, true, nil
}

func (as *AuthService) ResetForgetPassword(userinfo UserInfo, newpassword string) error {
	user, err := as.GetUserByGmail(userinfo.Email)
	if err != nil {
		return err
	}
	if err := as.SetUserPassword(&user, newpassword); err != nil {
		return err
	}
	return nil
}