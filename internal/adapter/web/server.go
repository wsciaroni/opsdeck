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
	orgHandler *handler.OrgHandler,
	commentHandler *handler.CommentHandler,
	authMW *authMiddleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth Routes
	r.Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)
	r.Post("/auth/logout", authHandler.Logout)

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Method(http.MethodGet, "/health", NewHealthHandler(db))

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.Protect)
			r.Get("/me", authHandler.Me)

			// Admin Routes
			r.Get("/admin/export/tickets", ticketHandler.ExportTickets)

			r.Post("/tickets", ticketHandler.CreateTicket)
			r.Get("/tickets", ticketHandler.ListTickets)
			r.Get("/tickets/{ticketID}", ticketHandler.GetTicket)
			r.Patch("/tickets/{ticketID}", ticketHandler.UpdateTicket)

			// Comments
			r.Post("/tickets/{ticketID}/comments", commentHandler.Create)
			r.Get("/tickets/{ticketID}/comments", commentHandler.List)

			r.Post("/organizations", orgHandler.CreateOrganization)
			r.Post("/organizations/{id}/members", orgHandler.AddMember)
			r.Get("/organizations/{id}/members", orgHandler.ListMembers)
			r.Delete("/organizations/{id}/members/{userID}", orgHandler.RemoveMember)

			r.Get("/organizations/{id}/share", orgHandler.GetShareSettings)
			r.Put("/organizations/{id}/share", orgHandler.UpdateShareSettings)
			r.Post("/organizations/{id}/share/regenerate", orgHandler.RegenerateShareToken)
		})
	})

	// Static Files (Frontend)
	// This must be last to act as a catch-all
	staticHandler := NewStaticHandler(staticFS)
	staticHandler.Register(r)

	return r
}
