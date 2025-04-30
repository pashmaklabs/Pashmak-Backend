package serializers_profile

type CurrentProfileResponse struct {
	FirstName  string
	LastName   string
	Email      string
	Avatar_url string
	Score      uint
}

type GetProfileByIDResponse struct {
	FirstName  string
	LastName   string
	Avatar_url string
	Score      uint
}

type UpdateUserProfileRequest struct {
	FirstName string	`json:"firstname" binding:"required,alpha"`
	LastName  string	`json:"lastname" binding:"required,alpha"`
}
