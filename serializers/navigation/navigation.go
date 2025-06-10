package serializers_navigation

type WaypointsResponse struct {
	Routes []struct {
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry,omitempty"` // Use omitempty to handle cases where geometry is not present
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
}