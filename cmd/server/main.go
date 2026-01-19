package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"log/slog"

	"github.com/wsciaroni/opsdeck/internal/adapter/auth/google"
	"github.com/wsciaroni/opsdeck/internal/adapter/storage"
	"github.com/wsciaroni/opsdeck/internal/adapter/storage/postgres"
	"github.com/wsciaroni/opsdeck/internal/adapter/web"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/service"
)

//go:embed all:dist
var webDist embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT/SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	// Database Configuration
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Run Database Migrations
	log.Println("Running database migrations...")
	if err := postgres.RunMigrations(ctx, databaseURL); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Connect to Database
	pool, err := postgres.ConnectPostgres(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	// Initialize River (Job Queue) - just ensuring it works, not running workers yet
	_, err = storage.InitRiver(ctx, pool)
	if err != nil {
		log.Fatalf("Failed to initialize River: %v", err)
	}
	log.Println("Initialized River client")

	// Prepare Static FS
	// dist is the root of the embedded FS, but we are inside cmd/server so it is relative to that?
	// Actually, the embed is relative to the file.
	// Since we copied web/dist to cmd/server/dist, the root is "dist".
	staticFS, err := fs.Sub(webDist, "dist")
	if err != nil {
		log.Fatalf("Failed to create sub FS for frontend: %v", err)
	}

	// Config
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if googleClientID == "" || googleClientSecret == "" || googleCallbackURL == "" {
		log.Fatal("Missing required environment variables: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Init Auth
	repo := postgres.NewUserRepository(pool)
	orgRepo := postgres.NewOrganizationRepository(pool)
	oidcProvider, err := google.NewGoogleProvider(ctx, googleClientID, googleClientSecret, googleCallbackURL)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}
	authService := service.NewAuthService(repo, orgRepo, oidcProvider, logger)
	authHandler := handler.NewAuthHandler(authService, logger)

	// Init Ticket
	ticketRepo := postgres.NewTicketRepository(pool)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService, orgRepo, logger)

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(repo, logger)

	// Setup Router
	router := web.NewRouter(pool, staticFS, authHandler, ticketHandler, authMiddleware)

	// Start Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")
}
