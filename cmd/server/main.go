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

	"github.com/wsciaroni/opsdeck/internal/adapter/storage"
	"github.com/wsciaroni/opsdeck/internal/adapter/web"
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

	// Connect to Database
	pool, err := storage.ConnectPostgres(ctx)
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

	// Setup Router
	router := web.NewRouter(pool, staticFS)

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
