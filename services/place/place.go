package services_place

import (
	"fmt"
	
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	sp "pashmak.com/pashmak/serializers/place"
)

type PlaceService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}

func NewPlaceService(db *gorm.DB, appconfig *bootstrap.AppConfig) *PlaceService {
	return &PlaceService{
		DB:        db,
		AppConfig: appconfig,
	}
}

func (ps *PlaceService) GetPlaceByID(id uint) (*sp.GetPlaceByIDResponse, error) {
	var results []sp.GetPlaceByIDResponse

	query := `
        SELECT 
            name,
            amenity,
            ST_Y(ST_Transform(way, 4326)) as latitude,
            ST_X(ST_Transform(way, 4326)) as longitude
        FROM planet_osm_point
        WHERE osm_id = ?`

	err := ps.DB.Raw(query, id).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no place found with ID %d", id)
	} else if len(results) > 1 {
		return nil, fmt.Errorf("multiple places found with ID %d", id)
	}

	return &results[0], nil
}

func (ps *PlaceService) SearchPlace(name string) ([]sp.GetPlaceByIDResponse, error) {
	var results []sp.GetPlaceByIDResponse

	query := `
		SELECT 
			name,
			amenity,
			ST_Y(ST_Transform(way, 4326)) as latitude,
			ST_X(ST_Transform(way, 4326)) as longitude
		FROM planet_osm_point
		WHERE name ILIKE ?
		LIMIT 10`

	err := ps.DB.Raw(query, "%"+name+"%").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
