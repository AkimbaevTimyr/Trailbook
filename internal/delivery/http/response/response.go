package response

import (
	"encoding/json"
	"net/http"
)

// Response is a unified API response envelope.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// JSON writes a raw JSON response with the given status code.
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success writes a 200 OK response wrapping data in the unified envelope.
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{Data: data})
}

// Created writes a 201 Created response wrapping data in the unified envelope.
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Response{Data: data})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes an error response with the given status code and message.
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, Response{Error: message})
}
