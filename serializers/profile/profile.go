package serializers_profile

import "time"

type CurrentProfileResponse struct {
	FirstName string
	LastName  string
	BirthDate time.Time //handle
	AboutMe   string
}
