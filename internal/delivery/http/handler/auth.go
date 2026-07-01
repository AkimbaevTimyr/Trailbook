package handler

import (
	"encoding/json"
	"net/http"
	"tracking-backend/internal/delivery/http/requests"

	"tracking-backend/internal/delivery/http/response"
	"tracking-backend/internal/usecase"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	uc usecase.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Login authenticates a user and returns an access token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokenResp, err := h.uc.Login(r.Context(), req)
	if err != nil {
		if isAuthError(err) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, tokenResp)
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case "email and password are required", "invalid email or password":
		return true
	}
	return false
}
