package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse es la estructura estándar para todas las respuestas de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON envía una respuesta JSON genérica
func JSON(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// Success envía una respuesta exitosa (200 OK)
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, true, message, data)
}

// Created envía una respuesta de recurso creado (201 Created)
func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, true, message, data)
}

// Error envía una respuesta de error estandarizada
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, false, message, nil)
}

// BadRequest envía un error 400
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// Unauthorized envía un error 401
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

// Forbidden envía un error 403
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

// NotFound envía un error 404
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalServerError envía un error 500
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
