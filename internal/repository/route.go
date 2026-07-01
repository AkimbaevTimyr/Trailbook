package repository

import (
	"context"
	"database/sql"
	"fmt"

	"tracking-backend/internal/domain"
)

// RouteRepository defines the interface for route data access.
type RouteRepository interface {
	GetAll(ctx context.Context) ([]domain.Route, error)
	GetAllByUserID(ctx context.Context, userID int64) ([]domain.Route, error)
	GetByID(ctx context.Context, id int64) (*domain.Route, error)
	Create(ctx context.Context, route *domain.Route) (*domain.Route, error)
	Update(ctx context.Context, route *domain.Route) (*domain.Route, error)
	Delete(ctx context.Context, id int64) error
}

type routeRepository struct {
	db *sql.DB
}

// NewRouteRepository creates a new RouteRepository.
func NewRouteRepository(db *sql.DB) RouteRepository {
	return &routeRepository{db: db}
}

func scanRoute(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Route, error) {
	var route domain.Route
	var description sql.NullString
	var regionID sql.NullInt64
	var createdBy sql.NullInt64
	var distanceKm sql.NullFloat64
	var elevationGain sql.NullInt64
	var durationHours sql.NullFloat64
	var season sql.NullString
	var trackGeoJSON sql.NullString
	var gpxFileURL sql.NullString
	var updatedAt sql.NullTime

	err := scanner.Scan(
		&route.ID,
		&route.Name,
		&description,
		&regionID,
		&createdBy,
		&route.RouteType,
		&route.Difficulty,
		&distanceKm,
		&elevationGain,
		&durationHours,
		&season,
		&trackGeoJSON,
		&gpxFileURL,
		&route.IsDeleted,
		&route.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return route, err
	}

	if description.Valid {
		route.Description = &description.String
	}
	if regionID.Valid {
		route.RegionID = &regionID.Int64
	}
	if createdBy.Valid {
		route.CreatedBy = &createdBy.Int64
	}
	if distanceKm.Valid {
		route.DistanceKm = &distanceKm.Float64
	}
	if elevationGain.Valid {
		eg := int(elevationGain.Int64)
		route.ElevationGain = &eg
	}
	if durationHours.Valid {
		route.DurationHours = &durationHours.Float64
	}
	if season.Valid {
		route.Season = &season.String
	}
	if trackGeoJSON.Valid {
		route.TrackGeoJSON = &trackGeoJSON.String
	}
	if gpxFileURL.Valid {
		route.GpxFileURL = &gpxFileURL.String
	}
	if updatedAt.Valid {
		route.UpdatedAt = &updatedAt.Time
	}

	return route, nil
}

func (r *routeRepository) GetAll(ctx context.Context) ([]domain.Route, error) {
	return r.queryRoutes(ctx, "is_deleted = FALSE", nil)
}

func (r *routeRepository) GetAllByUserID(ctx context.Context, userID int64) ([]domain.Route, error) {
	return r.queryRoutes(ctx, "is_deleted = FALSE AND created_by = $1", []interface{}{userID})
}

func (r *routeRepository) queryRoutes(ctx context.Context, where string, args []interface{}) ([]domain.Route, error) {
	query := `
		SELECT id, name, description, region_id, created_by, route_type, difficulty,
		       distance_km, elevation_gain, duration_hours, season, track_geojson,
		       gpx_file_url, is_deleted, created_at, updated_at
		FROM routes
		WHERE ` + where + `
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query routes: %w", err)
	}
	defer rows.Close()

	var routes []domain.Route
	for rows.Next() {
		route, err := scanRoute(rows)
		if err != nil {
			return nil, fmt.Errorf("scan route: %w", err)
		}
		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate routes: %w", err)
	}

	return routes, nil
}

func (r *routeRepository) GetByID(ctx context.Context, id int64) (*domain.Route, error) {
	route, err := scanRoute(r.db.QueryRowContext(ctx, `
		SELECT id, name, description, region_id, created_by, route_type, difficulty,
		       distance_km, elevation_gain, duration_hours, season, track_geojson,
		       gpx_file_url, is_deleted, created_at, updated_at
		FROM routes
		WHERE id = $1 AND is_deleted = FALSE
	`, id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query route by id: %w", err)
	}

	return &route, nil
}

func (r *routeRepository) Create(ctx context.Context, route *domain.Route) (*domain.Route, error) {
	query := `
		INSERT INTO routes (name, description, region_id, created_by, route_type, difficulty,
		                    distance_km, elevation_gain, duration_hours, season, track_geojson, gpx_file_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, description, region_id, created_by, route_type, difficulty,
		          distance_km, elevation_gain, duration_hours, season, track_geojson,
		          gpx_file_url, is_deleted, created_at, updated_at
	`

	created, err := scanRoute(r.db.QueryRowContext(
		ctx,
		query,
		route.Name,
		route.Description,
		route.RegionID,
		route.CreatedBy,
		route.RouteType,
		route.Difficulty,
		route.DistanceKm,
		route.ElevationGain,
		route.DurationHours,
		route.Season,
		route.TrackGeoJSON,
		route.GpxFileURL,
	))
	if err != nil {
		return nil, fmt.Errorf("create route: %w", err)
	}

	return &created, nil
}

func (r *routeRepository) Update(ctx context.Context, route *domain.Route) (*domain.Route, error) {
	query := `
		UPDATE routes
		SET name = $1,
		    description = $2,
		    region_id = $3,
		    route_type = $4,
		    difficulty = $5,
		    distance_km = $6,
		    elevation_gain = $7,
		    duration_hours = $8,
		    season = $9,
		    track_geojson = $10,
		    gpx_file_url = $11,
		    updated_at = NOW()
		WHERE id = $12 AND is_deleted = FALSE
		RETURNING id, name, description, region_id, created_by, route_type, difficulty,
		          distance_km, elevation_gain, duration_hours, season, track_geojson,
		          gpx_file_url, is_deleted, created_at, updated_at
	`

	updated, err := scanRoute(r.db.QueryRowContext(
		ctx,
		query,
		route.Name,
		route.Description,
		route.RegionID,
		route.RouteType,
		route.Difficulty,
		route.DistanceKm,
		route.ElevationGain,
		route.DurationHours,
		route.Season,
		route.TrackGeoJSON,
		route.GpxFileURL,
		route.ID,
	))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update route: %w", err)
	}

	return &updated, nil
}

func (r *routeRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE routes
		SET is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`, id)
	if err != nil {
		return fmt.Errorf("delete route: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
