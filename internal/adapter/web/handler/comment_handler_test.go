package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// MockCommentService is already defined in public_view_handler_test.go
// MockTicketService, MockOrgRepo, MockUserRepo are defined in ticket_handler_test.go

func TestListComments_Optimized(t *testing.T) {
	mockCommentService := new(MockCommentService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	mockTicketService := new(MockTicketService)

	h := handler.NewCommentHandler(mockCommentService, mockTicketService, mockUserRepo, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Get("/tickets/{ticketID}/comments", h.List)

	ticketID := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	commentAuthorID := uuid.New()

	user := &domain.User{ID: userID, Role: domain.RoleStaff}
	ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}

	comments := []domain.Comment{
		{
			ID:        uuid.New(),
			TicketID:  ticketID,
			UserID:    commentAuthorID,
			Body:      "Test Comment",
			CreatedAt: time.Now(),
		},
	}

	// 1. GetTicket
	mockTicketService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)

	// 2. CheckOrgAccess (Expect Optimization: GetMemberRole)
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, userID).Return("member", nil)

	// 3. ListComments
	mockCommentService.On("ListComments", mock.Anything, ticketID, true).Return(comments, nil)

	// 4. Fetch Users (Expect Optimization: GetByIDs)
	mockUserRepo.On("GetByIDs", mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
		return len(ids) == 1 && ids[0] == commentAuthorID
	})).Return([]domain.User{
		{ID: commentAuthorID, Name: "Author"},
	}, nil)

	req := httptest.NewRequest("GET", "/tickets/"+ticketID.String()+"/comments", nil)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Verification
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []handler.CommentResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Author", resp[0].User.Name)
}
