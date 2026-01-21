package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type AuthHandler struct {
	service port.AuthService
	orgRepo port.OrganizationRepository
	logger  *slog.Logger
	secret  []byte
}

func NewAuthHandler(service port.AuthService, orgRepo port.OrganizationRepository, logger *slog.Logger, secret string) *AuthHandler {
	return &AuthHandler{
		service: service,
		orgRepo: orgRepo,
		logger:  logger,
		secret:  []byte(secret),
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

	secure := true
	if os.Getenv("APP_ENV") == "development" {
		secure = false
	}

	signedSessionID := middleware.SignSessionID(sessionID, h.secret)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    signedSessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Set a reasonable expiration
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	secure := true
	if os.Getenv("APP_ENV") == "development" {
		secure = false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}

type MeResponse struct {
	User          *domain.User            `json:"user"`
	Organizations []domain.UserMembership `json:"organizations"`
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	organizations, err := h.orgRepo.ListByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list user organizations", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := MeResponse{
		User:          user,
		Organizations: organizations,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to write response", "error", err)
	}
}
