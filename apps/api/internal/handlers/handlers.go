package handlers

import (
	"encoding/json"
	"net/http"
	"velvet-archive-api/internal/db"
)

// BaseHandler holds our database queries wrapper (Dependency Injection)
type BaseHandler struct {
	DB *db.Queries
}

// NewBaseHandler constructs a new handler instance with our sqlc store
func NewBaseHandler(queries *db.Queries) *BaseHandler {
	return &BaseHandler{
		DB: queries,
	}
}

// respondWithJSON is a helper to write consistent JSON responses
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			logError(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

// respondWithError sends a structured error message back to the client
func respondWithError(w http.ResponseWriter, status int, msg string) {
	respondWithJSON(w, status, map[string]string{"error": msg})
}

func logError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
}
