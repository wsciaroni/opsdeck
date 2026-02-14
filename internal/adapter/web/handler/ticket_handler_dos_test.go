package handler_test

import (
	"context"
	"io"
	"mime/multipart"
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

type SpyReader struct {
	io.Reader
	TotalRead int
}

func (r *SpyReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.TotalRead += n
	return n, err
}

func TestCreatePublicTicket_DoSPrevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
	r := chi.NewRouter()
	r.Post("/public/tickets", h.CreatePublicTicket)

	t.Run("Rejects invalid token in URL WITHOUT reading body", func(t *testing.T) {
		// Setup large body (but we assert it is NOT read)
		// We use a pipe to simulate a stream that we can monitor, but SpyReader on a standard reader is easier.
		// If we use a huge reader, and the handler reads it all, the test might be slow.
		// So we use a moderate size, but track reads.
		// Actually, ParseMultipartForm reads the whole body if it fits in memory/disk.

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer writer.Close()
			defer pw.Close()
			// Write some fields
			_ = writer.WriteField("title", "Test Ticket")
			// We don't need to write a lot, just enough to prove it started reading.
			// But to simulate DoS, the point is we want to stop BEFORE reading.
		}()

		spy := &SpyReader{Reader: pr}

		// Token in URL is invalid
		invalidToken := "invalid-token"

		mockOrgRepo.On("GetByShareToken", mock.Anything, invalidToken).Return(nil, nil).Once()

		req := httptest.NewRequest("POST", "/public/tickets?token="+invalidToken, spy)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Expect 403 Forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Expect almost no body read (maybe some initial buffering by http server, but definitely not full parse)
		// ParseMultipartForm typically reads the whole body.
		// If our fix works, it returns 403 BEFORE calling ParseMultipartForm.
		// So TotalRead should be minimal (header sniffing).
		// Note: httptest might behave slightly differently than real server regarding buffering,
		// but SpyReader wraps the Body which is read by the handler.

		// If the handler calls ParseMultipartForm, it will try to read until EOF or error.
		// Since our pipe writer closes, it would read everything.
		// If we assert TotalRead == 0, it might be too strict if there's some buffering?
		// Let's print it to see.
		t.Logf("Total bytes read: %d", spy.TotalRead)

		// In a successful DoS prevention, we shouldn't read the body at all because we check URL param first.
		assert.Equal(t, 0, spy.TotalRead, "Should not read body if token is invalid")
	})

	t.Run("Requires token in URL for Multipart", func(t *testing.T) {
		// If token is missing from URL, it should return 400 Bad Request WITHOUT reading body

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			writer.Close()
			pw.Close()
		}()

		spy := &SpyReader{Reader: pr}

		req := httptest.NewRequest("POST", "/public/tickets", spy) // No token in URL
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, 0, spy.TotalRead, "Should not read body if token missing from URL")
	})
}

func TestCreateTicket_DoSPrevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
	r := chi.NewRouter()
	r.Post("/tickets", h.CreateTicket)

	t.Run("Rejects invalid orgID in URL WITHOUT reading body", func(t *testing.T) {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			writer.Close()
			pw.Close()
		}()
		spy := &SpyReader{Reader: pr}

		orgID := uuid.New()
		user := &domain.User{ID: uuid.New()}

		// User is NOT member
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("", nil).Once()

		req := httptest.NewRequest("POST", "/tickets?organization_id="+orgID.String(), spy)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, 0, spy.TotalRead, "Should not read body if not member")
	})

	t.Run("Requires organization_id in URL for Multipart", func(t *testing.T) {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			writer.Close()
			pw.Close()
		}()
		spy := &SpyReader{Reader: pr}

		user := &domain.User{ID: uuid.New()}

		req := httptest.NewRequest("POST", "/tickets", spy) // No orgID
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, 0, spy.TotalRead, "Should not read body if orgID missing")
	})
}
