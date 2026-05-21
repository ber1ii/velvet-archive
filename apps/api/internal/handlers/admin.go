package handlers

import (
	"encoding/json"
	"net/http"
	"velvet-archive-api/internal/db"

	"github.com/go-chi/chi/v5"
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

	respondWithJSON(w, http.StatusCreated, newSeries)
}

// Structural shape for creating a new lore entry inside a series
type CreateEntryRequest struct {
	SeriesID string `json:"series_id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
}

// POST /api/v1/admin/entries - Create a new lore entry inside a series
func (bh *BaseHandler) AdminCreateEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON input body")
		return
	}

	// Validate critical parameters
	if req.SeriesID == "" || req.Title == "" || req.Category == "" {
		respondWithError(w, http.StatusBadRequest, "Series ID, Title, and Category are required parameters")
		return
	}

	// Unpack structural string ID into Postgres native UUID representation
	var seriesUUID pgtype.UUID
	if err := seriesUUID.Scan(req.SeriesID); err != nil || !seriesUUID.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid target series UUID format")
		return
	}

	// Handle Metadata field processing safely
	// If it's a JSON column, we want to serialize our string summary into a valid JSON string structure
	var metadataBytes []byte
	if req.Summary != "" {
		var err error
		metadataBytes, err = json.Marshal(map[string]string{"summary": req.Summary})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to build metadata structural block")
			return
		}
	} else {
		metadataBytes = []byte("{}")
	}

	// Persist the entry using sqlc, completely aligning with the types generated from the database schema
	newEntry, err := bh.DB.CreateLoreEntry(r.Context(), db.CreateLoreEntryParams{
		SeriesID: seriesUUID,
		Title:    req.Title,
		Category: req.Category,
		Content:  req.Content,
		Metadata: metadataBytes,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to write lore entry")
		return
	}

	respondWithJSON(w, http.StatusCreated, newEntry)
}

// PUT /api/v1/admin/links/{id}/reveal - Flip link visibility to public view
func (bh *BaseHandler) AdminRevealLink(w http.ResponseWriter, r *http.Request) {
	// Grab the targeted relationship link ID from the routing param
	idParam := chi.URLParam(r, "id")

	var linkUUID pgtype.UUID
	if err := linkUUID.Scan(idParam); err != nil || !linkUUID.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid link relationship UUID format")
		return
	}

	// Execute the mutation update statement
	updatedLink, err := bh.DB.RevealLoreLink(r.Context(), linkUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Lore relationship link reference not found")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedLink)
}
