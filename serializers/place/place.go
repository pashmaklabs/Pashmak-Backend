package serializers_place

type GetPlaceByIDResponse struct {
	Name      string  `json:"name"`
	Amenity   string  `json:"amenity"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ID        int64   `json:"id"`
}

type PlaceWithRatingResponse struct {
	Name      string  `json:"name"`
	Amenity   string  `json:"amenity"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    float64 `json:"rating"`
}
