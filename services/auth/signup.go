package services_auth

import(
	serializers_auth "pashmak.com/pashmak/serializers"
)

func (as *AuthService)SignUp(Email string) (bool, error) {
	user, err := as.GetUserByGmail(Email)
	if err != nil {
		return false, err
	}

	
}