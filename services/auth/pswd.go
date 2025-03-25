package services_auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	models_auth "pashmak.com/pashmak/models"
)

func SendResetPasswordMail(email string, userOTP string) error {
	err := SendMail(email, userOTP)
	return err
}

func (as *AuthService) SetUserPassword(user *models_auth.User, newpassword string) error {
	user.Password = newpassword
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
	if user.Password != newpassword {
		return nil, nil
	}
	return &user, nil
}

func (as *AuthService) LoginWithPassword(email string, newpassword string) (bool, error) {
	if user, err := as.CheckUserPassword(email, newpassword); err != nil {
		return false, err
	} else if user == nil {
		return false, nil
	} else {
		// jwt, err := ac.authService.GenerateJWT(user)
		// TODO: Generate jwt and set cookie
		return true, nil
	}
}

func (as *AuthService) ForgetPassword(email string) error {
	// user, err := as.GetUserByGmail(email)
	// if err != nil {
	// 	return err
	// }
	userOTP := GenerateOTP()
	ctx := context.Background()
	if err := as.RedisClient.Set(ctx, email, userOTP, 2*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}
	if err := SendResetPasswordMail(email, userOTP); err != nil {
		return err
	}

	return nil
}

func (as *AuthService) VerifyForgetPassword(email string, otp string) (bool, error) {
	ctx := context.Background()
	realOTP, err := as.RedisClient.Get(ctx, email).Result()
	if err != nil {
		return false, fmt.Errorf("failed to get OTP from redis: %w", err)
	}

	if realOTP != otp {
		return false, nil
	}

	return true, nil
}
