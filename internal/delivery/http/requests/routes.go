package requests

// CreateRouteRequest represents a request to create a new route.
type CreateRouteRequest struct {
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	RegionID      *int64   `json:"region_id,omitempty"`
	RouteType     string   `json:"route_type"`
	Difficulty    string   `json:"difficulty"`
	DistanceKm    *float64 `json:"distance_km,omitempty"`
	ElevationGain *int     `json:"elevation_gain,omitempty"`
	DurationHours *float64 `json:"duration_hours,omitempty"`
	Season        *string  `json:"season,omitempty"`
}

// UpdateRouteRequest represents a request to update an existing route.
type UpdateRouteRequest struct {
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	RegionID      *int64   `json:"region_id,omitempty"`
	RouteType     string   `json:"route_type"`
	Difficulty    string   `json:"difficulty"`
	DistanceKm    *float64 `json:"distance_km,omitempty"`
	ElevationGain *int     `json:"elevation_gain,omitempty"`
	DurationHours *float64 `json:"duration_hours,omitempty"`
	Season        *string  `json:"season,omitempty"`
}
