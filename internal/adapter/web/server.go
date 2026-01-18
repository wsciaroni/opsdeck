package web

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, staticFS fs.FS) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Method(http.MethodGet, "/health", NewHealthHandler(db))
	})

	// Static Files (Frontend)
	// This must be last to act as a catch-all
	staticHandler := NewStaticHandler(staticFS)
	staticHandler.Register(r)

	return r
}
