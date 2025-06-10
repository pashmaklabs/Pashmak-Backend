package placeOsmUtils

import (
	"errors"

	"gorm.io/gorm"
	models_place "pashmak.com/pashmak/models/place"
)

func ImportFromOSM(placeID uint, DB *gorm.DB) (*models_place.Place, error) {
	var place models_place.Place
	if err := DB.Where("id = ?", placeID).First(&place).Error; err != nil {
		var count int64
		err := DB.Raw(`
			SELECT COUNT(*)
			FROM planet_osm_point
			WHERE osm_id = ?
		`, placeID).Scan(&count).Error
		if err != nil {
			return nil, err
		}
		if count > 0 {
			// Fetch OSM data
			var osmData struct {
				Name      string
				Amenity   string
				Latitude  float64
				Longitude float64
			}
			err := DB.Raw(`
				SELECT name, amenity, ST_Y(ST_Transform(way, 4326)) as latitude, ST_X(ST_Transform(way, 4326)) as longitude
				FROM planet_osm_point
				WHERE osm_id = ?
			`, placeID).Scan(&osmData).Error
			if err != nil {
				return nil, err
			}
			res := DB.Create(&models_place.Place{
				ID:        placeID,
				IsOSM:     true,
				Name:      osmData.Name,
				Amenity:   osmData.Amenity,
				Latitude:  osmData.Latitude,
				Longitude: osmData.Longitude,
			})
			if res.Error != nil {
				return nil, res.Error
			}
			// Retrieve the newly created place
			if err := DB.Where("id = ?", placeID).First(&place).Error; err != nil {
				return nil, err
			}
		} else if count == 0 {
			return nil, errors.New("place not found")
		}
	}
	return &place, nil
}