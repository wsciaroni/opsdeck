package web

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	authMiddleware "github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
)

func NewRouter(
	db *pgxpool.Pool,
	staticFS fs.FS,
	authHandler *handler.AuthHandler,
	ticketHandler *handler.TicketHandler,
	authMW *authMiddleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth Routes
	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Method(http.MethodGet, "/health", NewHealthHandler(db))
		r.Get("/me", authHandler.Me)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.Protect)
			r.Post("/tickets", ticketHandler.CreateTicket)
			r.Get("/tickets", ticketHandler.ListTickets)
			r.Patch("/tickets/{ticketID}", ticketHandler.UpdateTicket)
		})
	})

	// Static Files (Frontend)
	// This must be last to act as a catch-all
	staticHandler := NewStaticHandler(staticFS)
	staticHandler.Register(r)

	return r
}
