package web

import (
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
)

type StaticHandler struct {
	fs fs.FS
}

func NewStaticHandler(fs fs.FS) *StaticHandler {
	return &StaticHandler{fs: fs}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (h *StaticHandler) Register(r chi.Router) {
	if _, err := fs.Stat(h.fs, "index.html"); err != nil {
		// TODO: handle this specific error case later
		// If index.html is missing, it might be in a subdirectory (e.g. dist)
		// but here we expect the root of h.fs to contain index.html
		// For now we will assume the passed fs is correct.
		// In a real scenario we might want to panic or log error if the embedded fs is empty.
		log.Printf("warning: index.html not found in the provided filesystem: %v", err)
	}

	fileServer := http.FileServer(http.FS(h.fs))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		pathParam := chi.URLParam(r, "*")

		// If the file exists, serve it
		if fileExists(h.fs, pathParam) {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Otherwise, serve index.html (SPA routing)
		// But only if it is not an API call (which should have been handled before this catch-all)
		content, err := fs.ReadFile(h.fs, "index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write(content); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})
}

func fileExists(fileSystem fs.FS, filePath string) bool {
	// Clean the path to prevent directory traversal
	filePath = path.Clean(filePath)

	// If path is root or empty, we consider it a hit for index.html implicitly handled by serve or above logic
	if filePath == "." || filePath == "/" {
		return true
	}

	f, err := fileSystem.Open(filePath)
	if err != nil {
		return false
	}
	defer func() {
		_ = f.Close()
	}()

	s, err := f.Stat()
	if err != nil {
		return false
	}

	// We only want to serve actual files, not directories (unless we want directory listing which we don't for SPA)
	if s.IsDir() {
		return false
	}

	return true
}
