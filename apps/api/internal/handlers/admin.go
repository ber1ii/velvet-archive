package handlers

import (
	"encoding/json"
	"net/http"
	"velvet-archive-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// Struct to map the incoming JSON body for creating a Series
type CreateSeriesRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	CoverColor  string `json:"cover_color"`
	Description string `json:"description"`
}

// POST /api/v1/admin/series - Add a new bookshelf entry
func (bh *BaseHandler) AdminCreateSeries(w http.ResponseWriter, r *http.Request) {
	var req CreateSeriesRequest

	// Decode incoming JSON payload into the struct definition
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON input body")
		return
	}

	// Validate required fields
	if req.Title == "" || req.Author == "" || req.CoverColor == "" {
		respondWithError(w, http.StatusBadRequest, "Title, Author, and Cover Color are required fields")
		return
	}

	// Execute via sqlc wrapper
	newSeries, err := bh.DB.CreateSeries(r.Context(), db.CreateSeriesParams{
		Title:       req.Title,
		Author:      req.Author,
		CoverColor:  req.CoverColor,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to execute database write profile")
		return
	}

	// Respond with the created series
	respondWithJSON(w, http.StatusCreated, newSeries)
}
