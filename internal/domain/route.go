package domain

import "time"

// Route represents a hiking/alpine/ski route.
type Route struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	RegionID      *int64     `json:"region_id,omitempty"`
	CreatedBy     *int64     `json:"created_by,omitempty"`
	RouteType     string     `json:"route_type"`
	Difficulty    string     `json:"difficulty"`
	DistanceKm    *float64   `json:"distance_km,omitempty"`
	ElevationGain *int       `json:"elevation_gain,omitempty"`
	DurationHours *float64   `json:"duration_hours,omitempty"`
	Season        *string    `json:"season,omitempty"`
	TrackGeoJSON  *string    `json:"track_geojson,omitempty"`
	GpxFileURL    *string    `json:"gpx_file_url,omitempty"`
	IsDeleted     bool       `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
