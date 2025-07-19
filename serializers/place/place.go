package serializers_place

import models_place "pashmak.com/pashmak/models/place"

type GetPlaceByIDResponse struct {
	OsmID     *uint    `json:"osm_id"`
	Name      string   `json:"name"`
	Amenity   *string  `json:"amenity"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	ID        int64    `json:"id"`
	ImageURLs []string `json:"image_urls"`
}

type PlaceWithRatingResponse struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Amenity    string   `json:"amenity"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Rating     float64  `json:"rating"`
	ImageURLs  []string `json:"image_urls"`
	PlaceLabel *models_place.PlaceLabel
}

type AddPlaceRequest struct {
	Name      string  `json:"name"`
	Amenity   string  `json:"amenity"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
