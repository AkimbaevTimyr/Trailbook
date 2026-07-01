package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tracking-backend/internal/delivery/http/middleware"
	"tracking-backend/internal/delivery/http/requests"
	"tracking-backend/internal/delivery/http/response"
	"tracking-backend/internal/usecase"
)

// RouteHandler handles HTTP requests for routes.
type RouteHandler struct {
	uc usecase.RouteUsecase
}

// NewRouteHandler creates a new RouteHandler.
func NewRouteHandler(uc usecase.RouteUsecase) *RouteHandler {
	return &RouteHandler{uc: uc}
}

// GetRoutes returns a list of route cards.
func (h *RouteHandler) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.uc.GetRoutes(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, routes)
}

func (h *RouteHandler) GetMyRoutes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routes, err := h.uc.GetRoutesByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, routes)
}

// GetRoute returns full information about a route.
func (h *RouteHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid route id")
		return
	}

	route, err := h.uc.GetRoute(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if route == nil {
		response.Error(w, http.StatusNotFound, "route not found")
		return
	}

	response.Success(w, route)
}

// CreateRoute creates a new route for the authenticated user.
func (h *RouteHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req requests.CreateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	route, err := h.uc.CreateRoute(r.Context(), userID, req)
	if err != nil {
		if isRouteValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, route)
}

// UpdateRoute updates an existing route.
func (h *RouteHandler) UpdateRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid route id")
		return
	}

	var req requests.UpdateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	route, err := h.uc.UpdateRoute(r.Context(), id, req)
	if err != nil {
		if isRouteValidationError(err) || err.Error() == "route not found" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, route)
}

// DeleteRoute soft-deletes a route.
func (h *RouteHandler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid route id")
		return
	}

	if err := h.uc.DeleteRoute(r.Context(), id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.Error(w, http.StatusNotFound, "route not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NoContent(w)
}

func isRouteValidationError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case "name is required", "invalid route_type", "invalid difficulty":
		return true
	}
	return false
}
