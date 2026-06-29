package domain

import "time"

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	RegionID  *int64    `json:"region_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Region represents a geographic region.
type Region struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
