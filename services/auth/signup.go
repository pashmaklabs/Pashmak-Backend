package services_auth

import (
	"fmt"

	serializers_auth "pashmak.com/pashmak/serializers/auth"
)

func (as *AuthService) SignUp(userinfo UserInfo, Payload serializers_auth.SignUpRequest) error {
	user, err := as.GetUserByGmail(userinfo.Email)
	if err != nil {
		return fmt.Errorf("failed to get user by gmail: %w", err)
	}
	user.FirstName = Payload.FirstName
	user.LastName = Payload.LastName
	if Payload.Password != Payload.PasswordConfirm {
		return ErrAuth.ErrPasswordMismatch
	}

	hashedpass, err := as.HashPassword(Payload.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedpass
	if err := as.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save new password in database: %w", err)
	}
	return nil
}
