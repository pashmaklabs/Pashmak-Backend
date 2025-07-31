package serializers_place

type GetPlaceByIDResponse struct {
	OsmID     *uint    `json:"osm_id"`
	Name      string   `json:"name"`
	Amenity   *string  `json:"amenity"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	ID        string   `json:"id"`
	ImageURLs []string `json:"image_urls"`
}

type AddPlaceRequest struct {
	Name      string  `json:"name"`
	Amenity   string  `json:"amenity"`
	Latitude  float64 `json:"latitude" binding:"latitude"`
	Longitude float64 `json:"longitude" binding:"longitude"`
}
type SavedLocationResponse struct {
	ID           int64 `json:"id"`
	PlaceLabelID int64 `json:"place_label_id"`
}

type PlaceWithRatingResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Amenity       string                 `json:"amenity"`
	Latitude      float64                `json:"latitude"`
	Longitude     float64                `json:"longitude"`
	Rating        float64                `json:"rating"`
	ImageURLs     []string               `json:"image_urls"`
	SavedLocation *SavedLocationResponse `json:"saved_location"`
}
