package handler

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type AuthHandler struct {
	service port.AuthService
	logger  *slog.Logger
}

func NewAuthHandler(service port.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	url := h.service.GetLoginURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != "state-random-string" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	user, err := h.service.LoginFromProvider(r.Context(), code)
	if err != nil {
		h.logger.Error("failed to login from provider", "error", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	sessionID, err := h.service.CreateSession(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to create session", "error", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	secure := false
	if os.Getenv("APP_ENV") == "production" {
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Set a reasonable expiration
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "guest"}`))
}
