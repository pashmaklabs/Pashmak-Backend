package services_auth

import (
	"errors"
	"log"

	serializers_auth "pashmak.com/pashmak/serializers/auth"
)

func (as *AuthService)SignUp(userinfo UserInfo, Payload serializers_auth.SignUpRequest) error {
	user, err := as.GetUserByGmail(userinfo.Email)
	log.Println("user", userinfo.Email)
	if err != nil {
		return err
	}
	// FIXME: Error in query of GetUserByGmail(email is nil)
	user.FirstName = Payload.FirstName
	user.LastName = Payload.LastName
	if Payload.Password != Payload.PasswordConfirm {
		return errors.New("passwords do not match")
	}

	hashedpass, err := as.HashPassword(Payload.Password)
	if err != nil {
		return err
	}
	user.Password = hashedpass
	if err := as.DB.Save(&user).Error; err != nil {
		return err 
	}
	return nil
}