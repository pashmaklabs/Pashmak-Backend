package authentication

type AuthService struct{}

func NewAuthService() *AuthService{
	return &AuthService{}
}


func (as *AuthService)ValidateUser(email string){
	
}