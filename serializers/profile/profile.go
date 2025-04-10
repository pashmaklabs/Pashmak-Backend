package serializers_profile

import "time"

type CurrentProfileResponse struct {
	FirstName string
	LastName  string
	Email     string
	BirthDate time.Time //handle
	AboutMe   string
}

type GetProfileByIDResponse struct {
	FirstName string
	LastName  string
	AboutMe   string
}

