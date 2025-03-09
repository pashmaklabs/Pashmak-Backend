package authentication

type AuthService struct{}

func NewAuthService() *AuthService{
	return &AuthService{}
}

func CheckExistance(email string) bool{
	
	return false
}

func (as *AuthService)ValidateUser(email string){
	
}