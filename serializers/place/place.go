package serializers_place

type GetPlaceByIDResponse struct {
	OsmID     uint     `json:"osm_id"`
	Name      string   `json:"name"`
	Amenity   string   `json:"amenity"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	ID        int64    `json:"id"`
	ImageURLs []string `json:"image_urls"`
}

type PlaceWithRatingResponse struct {
	Name      string   `json:"name"`
	Amenity   string   `json:"amenity"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Rating    float64  `json:"rating"`
	ImageURLs []string `json:"image_urls"`
}
