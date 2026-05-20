package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// GET /api/v1/entries/{id} - Single lore entry details
func (bh *BaseHandler) GetLoreEntryDetails(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var dbUUID pgtype.UUID
	if err := dbUUID.Scan(idParam); err != nil || !dbUUID.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid lore entry UUID format")
		return
	}

	// Fetch the specific lore block
	entry, err := bh.DB.GetLoreEntry(r.Context(), dbUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Lore entry not found")
		return
	}
	respondWithJSON(w, http.StatusOK, entry)
}

// GET /api/v1/entries/{id}/links - All REVEALED relationships pointing out from this entity
func (bh *BaseHandler) GetEntryLinks(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var dbUUID pgtype.UUID
	if err := dbUUID.Scan(idParam); err != nil || !dbUUID.Valid {
		respondWithError(w, http.StatusBadRequest, "Invalid entry UUID format")
		return
	}

	// This runs query that forces ll.is_revealed = TRUE
	links, err := bh.DB.GetRevealedLinksForEntry(r.Context(), dbUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch social links")
		return
	}

	respondWithJSON(w, http.StatusOK, links)
}
