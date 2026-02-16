package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

func TestCreatePublicTicket_FilenameSanitization(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		expectedFilename string
	}{
		{
			name:             "Path Traversal",
			filename:         "../../etc/passwd",
			expectedFilename: "passwd",
		},
		{
			name:             "Absolute Path",
			filename:         "/var/www/html/shell.php",
			expectedFilename: "shell.php",
		},
		{
			name:             "Header Injection Quotes",
			filename:         "test.txt\"; filename=\"evil.exe",
			expectedFilename: "test.txt__ filename=_evil.exe", // Quotes replaced
		},
		{
			name:             "Header Injection Semicolon",
			filename:         "test.txt; filename=evil.exe",
			expectedFilename: "test.txt_ filename=evil.exe", // Semicolon replaced
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTicketService)
			mockOrgRepo := new(MockOrgRepo)
			mockUserRepo := new(MockUserRepo)
			h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
			r := chi.NewRouter()
			r.Post("/public/tickets", h.CreatePublicTicket)

			token := "valid-token"
			orgID := uuid.New()
			org := &domain.Organization{
				ID:               orgID,
				ShareLinkEnabled: true,
				ShareLinkToken:   &token,
			}
			userID := uuid.New()
			user := &domain.User{
				ID:    userID,
				Email: "test@example.com",
				Role:  domain.RolePublic,
			}
			ticket := &domain.Ticket{
				ID:    uuid.New(),
				Title: "Security Test",
			}

			mockOrgRepo.On("GetByShareToken", mock.Anything, token).Return(org, nil)
			mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

			// Capture the filename passed to service
			var capturedFilename string
			mockService.On("CreateTicket", mock.Anything, mock.MatchedBy(func(cmd port.CreateTicketCmd) bool {
				if len(cmd.Files) > 0 {
					capturedFilename = cmd.Files[0].Filename
				}
				return true
			})).Return(ticket, nil)

			// Create multipart body
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			_ = writer.WriteField("token", token)
			_ = writer.WriteField("title", "Security Test")
			_ = writer.WriteField("description", "Desc")
			_ = writer.WriteField("name", "User")
			_ = writer.WriteField("email", "test@example.com")
			_ = writer.WriteField("priority_id", "medium")

			part, _ := writer.CreateFormFile("files", tc.filename)
			_, _ = part.Write([]byte("content"))
			_ = writer.Close()

			req := httptest.NewRequest("POST", "/public/tickets?token="+token, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			assert.Equal(t, tc.expectedFilename, capturedFilename, "Filename should be sanitized")
		})
	}
}

func TestCreateTicket_FilenameSanitization(t *testing.T) {
	// Similar test for authenticated endpoint
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
	r := chi.NewRouter()
	r.Post("/tickets", h.CreateTicket)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

	filename := "../../etc/passwd"
	expected := "passwd"

	mockService.On("CreateTicket", mock.Anything, mock.MatchedBy(func(cmd port.CreateTicketCmd) bool {
		return len(cmd.Files) > 0 && cmd.Files[0].Filename == expected
	})).Return(&domain.Ticket{ID: uuid.New()}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("organization_id", orgID.String())
	_ = writer.WriteField("title", "Title")
	_ = writer.WriteField("description", "Desc")
	_ = writer.WriteField("priority_id", "medium")

	part, _ := writer.CreateFormFile("files", filename)
	_, _ = part.Write([]byte("content"))
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/tickets?organization_id="+orgID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePublicTicket_DoS_Prevention(t *testing.T) {
	// Test that invalid token in query param returns 403 FAST (no body read implied)
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
	r := chi.NewRouter()
	r.Post("/public/tickets", h.CreatePublicTicket)

	// Mock OrgRepo to return error or nil for invalid token
	mockOrgRepo.On("GetByShareToken", mock.Anything, "invalid-token").Return(nil, nil)

	// Create a large body to simulate DoS attempt (though we won't actually send 32MB in test for speed)
	// But the handler should reject before reading it.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("token", "invalid-token")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/public/tickets?token=invalid-token", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockOrgRepo.AssertExpectations(t)
}

func TestCreateTicket_DoS_Prevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
	r := chi.NewRouter()
	r.Post("/tickets", h.CreateTicket)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()

	// User not member of org
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("organization_id", orgID.String())
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/tickets?organization_id="+orgID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockOrgRepo.AssertExpectations(t)
}

func TestCreateTicket_Mismatch_Prevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
	r := chi.NewRouter()
	r.Post("/tickets", h.CreateTicket)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()
	otherOrgID := uuid.New()

	// User is member of orgID (in query param)
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("organization_id", otherOrgID.String()) // Mismatch!
	_ = writer.WriteField("title", "Title")
	_ = writer.WriteField("description", "Desc")
	_ = writer.WriteField("priority_id", "medium")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/tickets?organization_id="+orgID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail due to mismatch
}

func TestCreatePublicTicket_Mismatch_Prevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
	r := chi.NewRouter()
	r.Post("/public/tickets", h.CreatePublicTicket)

	token := "valid-token"
	otherToken := "other-token"
	org := &domain.Organization{
		ID:               uuid.New(),
		ShareLinkEnabled: true,
		ShareLinkToken:   &token,
	}

	mockOrgRepo.On("GetByShareToken", mock.Anything, token).Return(org, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("token", otherToken) // Mismatch
	_ = writer.WriteField("title", "Title")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/public/tickets?token="+token, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTicketHandler_InputLength(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
	r := chi.NewRouter()
	r.Post("/tickets", h.CreateTicket)
	r.Patch("/tickets/{ticketID}", h.UpdateTicket)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleStaff}
	orgID := uuid.New()
	ticketID := uuid.New()
	ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}

	t.Run("CreateTicket - Location Too Long", func(t *testing.T) {
		// Called twice: once for query param check, once for body check
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil).Maybe()

		// Mock success because currently it SUCCEEDS (we expect failure in assertion, but code proceeds)
		mockService.On("CreateTicket", mock.Anything, mock.Anything).Return(ticket, nil).Maybe()

		reqBody := map[string]interface{}{
			"organization_id": orgID,
			"title":           "Valid Title",
			"description":     "Valid Description",
			"location":        strings.Repeat("A", 201), // Limit 200
			"priority_id":     "medium",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/tickets?organization_id="+orgID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Location")
	})

	t.Run("UpdateTicket - Location Too Long", func(t *testing.T) {
		mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil).Maybe()
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil).Maybe()

		// Mock success because currently it SUCCEEDS
		mockService.On("UpdateTicket", mock.Anything, ticketID, mock.Anything).Return(ticket, nil).Maybe()

		reqBody := map[string]interface{}{
			"location": strings.Repeat("A", 201),
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Location")
	})

	t.Run("UpdateTicket - AssigneeCustomName Too Long", func(t *testing.T) {
		mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil).Maybe()
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil).Maybe()

		// Mock success because currently it SUCCEEDS
		mockService.On("UpdateTicket", mock.Anything, ticketID, mock.Anything).Return(ticket, nil).Maybe()

		reqBody := map[string]interface{}{
			"assignee_custom_name": strings.Repeat("A", 101), // Limit 100
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Assignee Name")
	})
}
