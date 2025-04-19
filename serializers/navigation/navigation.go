package serializers_navigation

type WaypointsResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
			Type        string      `json:"type"`
		} `json:"geometry,omitempty"` // Use omitempty to handle cases where geometry is not present
		Legs []struct {
			Steps    []interface{} `json:"steps"`
			Summary  string        `json:"summary"`
			Weight   float64       `json:"weight"`
			Duration float64       `json:"duration"`
			Distance float64       `json:"distance"`
		} `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
}