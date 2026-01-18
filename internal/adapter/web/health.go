package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "connected"

	// Check database connection
	if err := h.db.Ping(r.Context()); err != nil {
		status = "error"
		dbStatus = "disconnected: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := map[string]string{
		"status":   status,
		"database": dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// We can't write a new HTTP status code here because the header is already sent,
		// but we should log it so we know something is wrong.
		log.Printf("failed to encode response: %v", err)
	}
}
