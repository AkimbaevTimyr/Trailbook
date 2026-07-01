package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tracking-backend/internal/delivery/http/requests"

	"tracking-backend/internal/delivery/http/response"
	"tracking-backend/internal/usecase"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetUsers(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.uc.GetUser(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if user == nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.uc.Create(r.Context(), req)
	if err != nil {
		if isValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, user)
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case "name is required", "email is required", "password must be at least 6 characters":
		return true
	}
	return false
}
