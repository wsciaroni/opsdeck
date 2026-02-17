package web

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	appMiddleware "github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
)

func NewRouter(
	db *pgxpool.Pool,
	staticFS fs.FS,
	authHandler *handler.AuthHandler,
	ticketHandler *handler.TicketHandler,
	orgHandler *handler.OrgHandler,
	commentHandler *handler.CommentHandler,
	publicViewHandler *handler.PublicViewHandler,
	scheduledTaskHandler *handler.ScheduledTaskHandler,
	authMW *appMiddleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.SecurityHeaders)

	// Auth Routes
	// Rate limit login initiation to 20 requests per minute per IP to prevent abuse/DoS
	loginRL := appMiddleware.NewRateLimiter(20, time.Minute)
	r.With(loginRL.Limit).Get("/auth/login", authHandler.Login)
	r.Get("/auth/callback", authHandler.Callback)
	r.Post("/auth/logout", authHandler.Logout)

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Method(http.MethodGet, "/health", NewHealthHandler(db))

		// Rate limit public ticket creation to 10 requests per minute per IP
		rl := appMiddleware.NewRateLimiter(10, time.Minute)
		r.With(rl.Limit).Post("/public/tickets", ticketHandler.CreatePublicTicket)

		r.Route("/public/view/{token}", func(r chi.Router) {
			r.Get("/organization", publicViewHandler.GetOrganization)
			r.Get("/tickets", publicViewHandler.ListTickets)
			r.Get("/tickets/{ticketID}", publicViewHandler.GetTicket)
			r.Get("/tickets/{ticketID}/comments", publicViewHandler.ListComments)
		})

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.Protect)
			r.Get("/me", authHandler.Me)

			// Admin Routes
			r.Get("/admin/export/tickets", ticketHandler.ExportTickets)

			// Rate limit ticket creation/updates to 20 requests per minute per IP
			// This mitigates DoS risks from authenticated users (e.g. large file uploads)
			ticketRL := appMiddleware.NewRateLimiter(20, time.Minute)
			r.With(ticketRL.Limit).Post("/tickets", ticketHandler.CreateTicket)
			r.Get("/tickets", ticketHandler.ListTickets)
			r.Get("/tickets/{ticketID}", ticketHandler.GetTicket)
			r.Get("/tickets/{ticketID}/files/{fileID}", ticketHandler.GetTicketFile)
			r.With(ticketRL.Limit).Patch("/tickets/{ticketID}", ticketHandler.UpdateTicket)

			// Comments
			r.With(ticketRL.Limit).Post("/tickets/{ticketID}/comments", commentHandler.Create)
			r.Get("/tickets/{ticketID}/comments", commentHandler.List)

			// Scheduled Tasks
			r.Get("/scheduled-tasks", scheduledTaskHandler.List)
			r.With(ticketRL.Limit).Post("/scheduled-tasks", scheduledTaskHandler.Create)
			r.Patch("/scheduled-tasks/{id}", scheduledTaskHandler.Update)
			r.Delete("/scheduled-tasks/{id}", scheduledTaskHandler.Delete)

			r.Post("/organizations", orgHandler.CreateOrganization)
			r.Post("/organizations/{id}/members", orgHandler.AddMember)
			r.Get("/organizations/{id}/members", orgHandler.ListMembers)
			r.Delete("/organizations/{id}/members/{userID}", orgHandler.RemoveMember)
			r.Put("/organizations/{id}/members/{userID}/role", orgHandler.UpdateMemberRole)

			r.Get("/organizations/{id}/share", orgHandler.GetShareSettings)
			r.Put("/organizations/{id}/share", orgHandler.UpdateShareSettings)
			r.Post("/organizations/{id}/share/regenerate", orgHandler.RegenerateShareToken)

			r.Get("/organizations/{id}/public-view", orgHandler.GetPublicViewSettings)
			r.Put("/organizations/{id}/public-view", orgHandler.UpdatePublicViewSettings)
			r.Post("/organizations/{id}/public-view/regenerate", orgHandler.RegeneratePublicViewToken)
		})
	})

	// Static Files (Frontend)
	// This must be last to act as a catch-all
	staticHandler := NewStaticHandler(staticFS)
	staticHandler.Register(r)

	return r
}
