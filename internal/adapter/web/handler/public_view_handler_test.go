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
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// MockCommentService (redefined here as it wasn't in ticket_handler_test.go)
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) CreateComment(ctx context.Context, cmd port.CreateCommentCmd) (*domain.Comment, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) ListComments(ctx context.Context, ticketID uuid.UUID, includeSensitive bool) ([]domain.Comment, error) {
	args := m.Called(ctx, ticketID, includeSensitive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func TestListComments_Hydration(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockTicketService := new(MockTicketService)
	mockCommentService := new(MockCommentService)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewPublicViewHandler(mockOrgRepo, mockTicketService, mockCommentService, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Get("/public/view/{token}/tickets/{ticketID}/comments", h.ListComments)

	token := "valid-token"
	orgID := uuid.New()
	org := &domain.Organization{
		ID:                orgID,
		PublicViewEnabled: true,
	}

	ticketID := uuid.New()
	ticket := &domain.Ticket{
		ID:             ticketID,
		OrganizationID: orgID,
		Sensitive:      false,
	}

	user1ID := uuid.New()
	user2ID := uuid.New()

	comments := []domain.Comment{
		{ID: uuid.New(), UserID: user1ID, Body: "Comment 1", CreatedAt: time.Now()},
		{ID: uuid.New(), UserID: user2ID, Body: "Comment 2", CreatedAt: time.Now()},
		{ID: uuid.New(), UserID: user1ID, Body: "Comment 3", CreatedAt: time.Now()},
	}

	mockOrgRepo.On("GetByPublicViewToken", mock.Anything, token).Return(org, nil)
	mockTicketService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)
	mockCommentService.On("ListComments", mock.Anything, ticketID, false).Return(comments, nil)

	// Optimized expectation: GetByIDs called once with all user IDs
	mockUserRepo.On("GetByIDs", mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
		// Verify all needed IDs are present
		hasUser1 := false
		hasUser2 := false
		for _, id := range ids {
			if id == user1ID {
				hasUser1 = true
			}
			if id == user2ID {
				hasUser2 = true
			}
		}
		return hasUser1 && hasUser2 && len(ids) == 2
	})).Return([]domain.User{
		{ID: user1ID, Name: "User 1"},
		{ID: user2ID, Name: "User 2"},
	}, nil)

	req := httptest.NewRequest("GET", "/public/view/"+token+"/tickets/"+ticketID.String()+"/comments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []handler.CommentResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 3)
	assert.Equal(t, "User 1", resp[0].User.Name)
	assert.Equal(t, "User 2", resp[1].User.Name)
	assert.Equal(t, "User 1", resp[2].User.Name)
}
