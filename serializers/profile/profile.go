package serializers_profile

type CurrentProfileResponse struct {
	FirstName string
	LastName  string
	Email     string
	Image_url string
	Score     uint
}

type GetProfileByIDResponse struct {
	FirstName string
	LastName  string
	Image_url string
	Score     uint
}
