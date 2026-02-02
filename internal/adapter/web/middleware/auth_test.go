package middleware_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.User, error) {
	return nil, nil
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func TestSignSessionID(t *testing.T) {
	secret := []byte("secret")
	id := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	signed := middleware.SignSessionID(id, expiresAt, secret)
	assert.NotEmpty(t, signed)
	assert.Contains(t, signed, id)

	parts := strings.Split(signed, ".")
	assert.Len(t, parts, 3)
}

func TestAuthMiddleware_Protect(t *testing.T) {
	secret := "secret"
	mockRepo := new(MockUserRepository)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mw := middleware.NewAuthMiddleware(mockRepo, logger, secret)
	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Valid Token", func(t *testing.T) {
		id := uuid.New()
		user := &domain.User{ID: id, Email: "test@example.com"}
		// Reset mock
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", mock.Anything, id).Return(user, nil)

		expiresAt := time.Now().Add(1 * time.Hour)
		token := middleware.SignSessionID(id.String(), expiresAt, []byte(secret))

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: token})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Expired Token", func(t *testing.T) {
		id := uuid.New()
		// No repo call expected because validation should fail first
		mockRepo.ExpectedCalls = nil

		expiresAt := time.Now().Add(-1 * time.Hour) // Past
		token := middleware.SignSessionID(id.String(), expiresAt, []byte(secret))

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: token})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Tampered Token", func(t *testing.T) {
		id := uuid.New()
		expiresAt := time.Now().Add(1 * time.Hour)
		token := middleware.SignSessionID(id.String(), expiresAt, []byte(secret))

		// Tamper signature
		token += "bad"

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: token})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Tampered Expiration", func(t *testing.T) {
		id := uuid.New()
		expiresAt := time.Now().Add(1 * time.Hour)
		// Properly signed
		token := middleware.SignSessionID(id.String(), expiresAt, []byte(secret))

		// Change exp to be valid for longer
		parts := strings.Split(token, ".")
		newExp := time.Now().Add(24 * time.Hour).Unix()
		fakeToken := parts[0] + "." + strconv.FormatInt(newExp, 10) + "." + parts[2]

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: fakeToken})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
