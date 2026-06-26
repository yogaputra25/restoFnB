package response

import (
	"encoding/json"
	"net/http"
	"time"
)

// Standard response envelope
type APIResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
	Data      any         `json:"data,omitempty"`
	Errors    []FieldErr  `json:"errors,omitempty"`
}

// Pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// Field validation error
type FieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ListPayload wraps data with pagination
type ListPayload struct {
	Items      any        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Success returns a 200 envelope with data.
func Success(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success:   true,
		Code:      http.StatusOK,
		Message:   "Success",
		Timestamp: timestamp(),
		Data:      data,
	})
}

// Created returns a 201 envelope with data.
func Created(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, APIResponse{
		Success:   true,
		Code:      http.StatusCreated,
		Message:   "Created",
		Timestamp: timestamp(),
		Data:      data,
	})
}

// List returns a 200 envelope with paginated data.
func List(w http.ResponseWriter, items any, page, limit, totalItems int) {
	totalPages := (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}
	writeJSON(w, http.StatusOK, APIResponse{
		Success:   true,
		Code:      http.StatusOK,
		Message:   "Success",
		Timestamp: timestamp(),
		Data: ListPayload{
			Items: items,
			Pagination: Pagination{
				Page:       page,
				Limit:      limit,
				TotalItems: totalItems,
				TotalPages: totalPages,
			},
		},
	})
}

// Error sends a simple error response.
func Error(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, APIResponse{
		Success:   false,
		Code:      status,
		Message:   msg,
		Timestamp: timestamp(),
	})
}

// ValidationError sends a 400 with per-field errors.
func ValidationError(w http.ResponseWriter, errs []FieldErr) {
	writeJSON(w, http.StatusBadRequest, APIResponse{
		Success:   false,
		Code:      http.StatusBadRequest,
		Message:   "Validation Error",
		Timestamp: timestamp(),
		Errors:    errs,
	})
}

// Unauthorized sends a 401.
func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "Unauthorized")
}

// Forbidden sends a 403.
func Forbidden(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "Forbidden")
}

// NotFound sends a 404.
func NotFound(w http.ResponseWriter) {
	Error(w, http.StatusNotFound, "Resource not found")
}

// ServerError sends a 500.
func ServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Internal Server Error")
}
