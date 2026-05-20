package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// GET /api/v1/series - List all series (bookshelves)
func (bh *BaseHandler) ListSeries(w http.ResponseWriter, r *http.Request) {
	seriesList, err := bh.DB.ListSeries(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve series archive")
		return
	}
	respondWithJSON(w, http.StatusOK, seriesList)
}

// GET /api/v1/series/{id} - Single series + its nested lore entries
func (bh *BaseHandler) GetSeriesDetails(w http.ResponseWriter, r *http.Request) {
	// Grab the ID parameter from the Chi router URL string
	idParam := chi.URLParam(r, "id")

	// Convert the string parameter directly into pgtype.UUID for sqlc compatibility
	var dbUUID pgtype.UUID
	err := dbUUID.Scan(idParam)
	if err != nil || !dbUUID.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Fetch the main series information using the correct type
	series, err := bh.DB.GetSeries(r.Context(), dbUUID)
	if err != nil {
		// Fixed the typo here from NZNotFound to StatusNotFound
		respondWithError(w, http.StatusNotFound, "Series not found")
		return
	}

	// Fetch all nested lore entries belonging to this series
	entries, err := bh.DB.GetLoreEntriesBySeries(r.Context(), dbUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load associated lore entries")
		return
	}

	// Combine the data structural response to send clean JSON
	response := map[string]interface{}{
		"series":  series,
		"entries": entries,
	}

	respondWithJSON(w, http.StatusOK, response)
}
