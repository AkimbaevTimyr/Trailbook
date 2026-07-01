package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"tracking-backend/internal/delivery/http/requests"
	"tracking-backend/internal/domain"
	"tracking-backend/internal/repository"
)

var (
	validRouteTypes   = map[string]bool{"hike": true, "alpine": true, "ski": true, "trail_run": true}
	validDifficulties = map[string]bool{"easy": true, "moderate": true, "hard": true, "extreme": true}
)

// RouteUsecase defines the interface for route-related business logic.
type RouteUsecase interface {
	GetRoutes(ctx context.Context) ([]domain.Route, error)
	GetRoutesByUserID(ctx context.Context, userID int64) ([]domain.Route, error)
	GetRoute(ctx context.Context, id int64) (*domain.Route, error)
	CreateRoute(ctx context.Context, userID int64, req requests.CreateRouteRequest) (*domain.Route, error)
	UpdateRoute(ctx context.Context, id int64, req requests.UpdateRouteRequest) (*domain.Route, error)
	DeleteRoute(ctx context.Context, id int64) error
}

type routeUsecase struct {
	repo repository.RouteRepository
}

// NewRouteUsecase creates a new RouteUsecase.
func NewRouteUsecase(repo repository.RouteRepository) RouteUsecase {
	return &routeUsecase{repo: repo}
}

func (uc *routeUsecase) GetRoutes(ctx context.Context) ([]domain.Route, error) {
	routes, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get routes: %w", err)
	}
	return routes, nil
}

func (uc *routeUsecase) GetRoutesByUserID(ctx context.Context, userID int64) ([]domain.Route, error) {
	routes, err := uc.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get routes by user: %w", err)
	}
	return routes, nil
}

func (uc *routeUsecase) GetRoute(ctx context.Context, id int64) (*domain.Route, error) {
	route, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get route: %w", err)
	}
	return route, nil
}

func (uc *routeUsecase) CreateRoute(ctx context.Context, userID int64, req requests.CreateRouteRequest) (*domain.Route, error) {
	if err := validateRouteRequest(req.Name, req.RouteType, req.Difficulty); err != nil {
		return nil, err
	}

	trackGeoJSON := "{}"
	gpxFileURL := ""

	route := &domain.Route{
		Name:          req.Name,
		Description:   req.Description,
		RegionID:      req.RegionID,
		CreatedBy:     &userID,
		RouteType:     req.RouteType,
		Difficulty:    req.Difficulty,
		DistanceKm:    req.DistanceKm,
		ElevationGain: req.ElevationGain,
		DurationHours: req.DurationHours,
		Season:        req.Season,
		TrackGeoJSON:  &trackGeoJSON,
		GpxFileURL:    &gpxFileURL,
	}

	created, err := uc.repo.Create(ctx, route)
	if err != nil {
		return nil, fmt.Errorf("create route: %w", err)
	}

	return created, nil
}

func (uc *routeUsecase) UpdateRoute(ctx context.Context, id int64, req requests.UpdateRouteRequest) (*domain.Route, error) {
	if err := validateRouteRequest(req.Name, req.RouteType, req.Difficulty); err != nil {
		return nil, err
	}

	trackGeoJSON := "{}"
	gpxFileURL := ""

	route := &domain.Route{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		RegionID:      req.RegionID,
		RouteType:     req.RouteType,
		Difficulty:    req.Difficulty,
		DistanceKm:    req.DistanceKm,
		ElevationGain: req.ElevationGain,
		DurationHours: req.DurationHours,
		Season:        req.Season,
		TrackGeoJSON:  &trackGeoJSON,
		GpxFileURL:    &gpxFileURL,
	}

	updated, err := uc.repo.Update(ctx, route)
	if err != nil {
		return nil, fmt.Errorf("update route: %w", err)
	}
	if updated == nil {
		return nil, errors.New("route not found")
	}

	return updated, nil
}

func (uc *routeUsecase) DeleteRoute(ctx context.Context, id int64) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete route: %w", err)
	}
	return nil
}

func validateRouteRequest(name, routeType, difficulty string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if !validRouteTypes[routeType] {
		return errors.New("invalid route_type")
	}
	if !validDifficulties[difficulty] {
		return errors.New("invalid difficulty")
	}
	return nil
}
