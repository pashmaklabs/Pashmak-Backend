package services_auth

import(
	serializers_auth "pashmak.com/pashmak/serializers"
	"errors"
)

func (as *AuthService)SignUp(Email string, Payload serializers_auth.SignUpRequest) error {
	user, err := as.GetUserByGmail(Email)
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