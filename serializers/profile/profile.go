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

type SavedLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserNote  string  `json:"user_note,omitempty"`
	Label     string  `json:"label"`
}

type SavedLocationRequest struct {
	PlaceID   *uint   `json:"place_id,omitempty"`
	Label     string  `json:"label" binding:"required,oneof=favorites to_go fun"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	UserNote  string  `json:"user_note,omitempty"`
}

type CreatePlaceLabelRequest struct {
	Name string
}

type GetSavedLocationsByPlaceLabelRequest struct {
	PlaceLabelID uint `uri:"place_label_id" form:"place_label_id"`
}

type UpdateSavedLocationRequest struct {
	ID           uint    `json:"id"`
	Name         *string `json:"name"`
	UserNote     *string `json:"user_note"`
	PlaceLabelID *uint   `json:"place_label_id"`
}

type DeleteSavedLocation struct {
	ID uint `uri:"id"`
}

type UpdateUserProfileRequest struct {
	FirstName string `json:"firstname" binding:"required,alpha"`
	LastName  string `json:"lastname" binding:"required,alpha"`
}

type CreateSavedLocationRequest struct {
	PlaceID      *uint   `json:"place_id,omitempty"`
	PlaceLabelID uint    `json:"place_label_id"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	UserNote     *string `json:"user_note,omitempty"`
	Name         *string `json:"name,omitempty"`
}

type PlaceLabelWithCountResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	SavedLocationsCount int64  `json:"saved_locations_count"`
}
