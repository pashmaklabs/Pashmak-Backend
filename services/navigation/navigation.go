package services_navigation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	serializers_navigation "pashmak.com/pashmak/serializers/navigation"
)

type NavigationService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}

func NewNavigationService(db *gorm.DB, appconfig *bootstrap.AppConfig) *NavigationService {
	return &NavigationService{
		DB:        db,
		AppConfig: appconfig,
	}
}


func (ns *NavigationService) FetchRoute(startLat, startLon, endLat, endLon string) (*serializers_navigation.WaypointsResponse, error) {
	osrmURL := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s,%s;%s,%s?overview=false&geometries=geojson",
		url.QueryEscape(startLon), url.QueryEscape(startLat),
		url.QueryEscape(endLon), url.QueryEscape(endLat))

	resp, err := http.Get(osrmURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM: %v", err)
	}
	defer resp.Body.Close()

	var waypoint serializers_navigation.WaypointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&waypoint); err != nil {
		return nil, fmt.Errorf("failed to decode OSRM response: %v", err)
	}

	RoutingURL := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%s,%s;%s,%s?overview=full&geometries=geojson",
		url.QueryEscape(fmt.Sprint(waypoint.Waypoints[0].Location[0])), url.QueryEscape(fmt.Sprint(waypoint.Waypoints[0].Location[1])),
		url.QueryEscape(fmt.Sprint(waypoint.Waypoints[1].Location[0])), url.QueryEscape(fmt.Sprint(waypoint.Waypoints[1].Location[1])))

		
	resp, err = http.Get(RoutingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM: %v", err)
	}

	var route serializers_navigation.WaypointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&route); err != nil {
		return nil, fmt.Errorf("failed to decode OSRM response: %v", err)
	}

	return &route, nil
}