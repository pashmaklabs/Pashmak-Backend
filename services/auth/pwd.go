package services_auth

import (
	"errors"

	models_auth "pashmak.com/pashmak/models"
)

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



func (as *AuthService)LoginWithPassword(email string, newpassword string) (bool, error) {
	if user, err := as.CheckUserPassword(email, newpassword); err != nil {
		return false, err
	} else if (user == nil) {
		return false, nil
	} else {
		as.SetUserPassword(user, newpassword)
		return true, nil
	}
}
