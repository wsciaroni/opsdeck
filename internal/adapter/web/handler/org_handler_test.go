package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

func TestGetShareSettings(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewOrgHandler(mockOrgRepo, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Get("/organizations/{id}/share", h.GetShareSettings)

	t.Run("Success", func(t *testing.T) {
		orgID := uuid.New()
		userID := uuid.New()
		user := &domain.User{ID: userID}

		token := "some-token"
		org := &domain.Organization{
			ID:               orgID,
			ShareLinkEnabled: true,
			ShareLinkToken:   &token,
		}

		mockOrgRepo.On("ListByUser", mock.Anything, userID).Return([]domain.UserMembership{
			{Organization: domain.Organization{ID: orgID}, Role: "member"},
		}, nil)
		mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(org, nil)

		req := httptest.NewRequest("GET", "/organizations/"+orgID.String()+"/share", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			ShareLinkEnabled bool    `json:"share_link_enabled"`
			ShareLinkToken   *string `json:"share_link_token"`
		}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.ShareLinkEnabled)
		assert.Equal(t, token, *resp.ShareLinkToken)
	})
}

func TestUpdateShareSettings(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewOrgHandler(mockOrgRepo, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Put("/organizations/{id}/share", h.UpdateShareSettings)

	t.Run("Success - Enable Share Link", func(t *testing.T) {
		orgID := uuid.New()
		userID := uuid.New()
		user := &domain.User{ID: userID}

		org := &domain.Organization{
			ID:               orgID,
			ShareLinkEnabled: false,
			ShareLinkToken:   nil,
		}

		mockOrgRepo.On("ListByUser", mock.Anything, userID).Return([]domain.UserMembership{
			{Organization: domain.Organization{ID: orgID}, Role: "admin"},
		}, nil)
		mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(org, nil)
		mockOrgRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *domain.Organization) bool {
			return o.ShareLinkEnabled == true && o.ShareLinkToken != nil
		})).Return(nil)

		body := map[string]bool{"enabled": true}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/organizations/"+orgID.String()+"/share", bytes.NewReader(bodyBytes))
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegenerateShareToken(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewOrgHandler(mockOrgRepo, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Post("/organizations/{id}/share/regenerate", h.RegenerateShareToken)

	t.Run("Success", func(t *testing.T) {
		orgID := uuid.New()
		userID := uuid.New()
		user := &domain.User{ID: userID}

		oldToken := "old-token"
		org := &domain.Organization{
			ID:               orgID,
			ShareLinkEnabled: true,
			ShareLinkToken:   &oldToken,
		}

		mockOrgRepo.On("ListByUser", mock.Anything, userID).Return([]domain.UserMembership{
			{Organization: domain.Organization{ID: orgID}, Role: "owner"},
		}, nil)
		mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(org, nil)
		mockOrgRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *domain.Organization) bool {
			return o.ShareLinkToken != nil && *o.ShareLinkToken != oldToken
		})).Return(nil)

		req := httptest.NewRequest("POST", "/organizations/"+orgID.String()+"/share/regenerate", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
