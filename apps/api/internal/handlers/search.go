package handlers

import "net/http"

// GET /api/v1/search?q=... - Deep full text search

func (bh *BaseHandler) SearchArchive(w http.ResponseWriter, r *http.Request) {
	// Parse URL Query parameters from request
	queryText := r.URL.Query().Get("q")
	if queryText == "" {
		respondWithError(w, http.StatusBadRequest, "Search query parameter 'q' cannot be empty")
		return
	}

	// Run our raw SQL index text search via sqlc
	results, err := bh.DB.SearchLoreEntries(r.Context(), queryText)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Search operation encountered an index error")
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}
