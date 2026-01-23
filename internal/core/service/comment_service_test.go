package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) ListByTicket(ctx context.Context, ticketID uuid.UUID, includeSensitive bool) ([]domain.Comment, error) {
	args := m.Called(ctx, ticketID, includeSensitive)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func TestCreateComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)
	ctx := context.Background()

	t.Run("Valid Comment", func(t *testing.T) {
		cmd := port.CreateCommentCmd{
			TicketID: uuid.New(),
			UserID:   uuid.New(),
			Body:     "Test comment",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Comment")).Return(nil).Once()

		comment, err := service.CreateComment(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, cmd.Body, comment.Body)
	})

	t.Run("Empty Body", func(t *testing.T) {
		cmd := port.CreateCommentCmd{
			TicketID: uuid.New(),
			UserID:   uuid.New(),
			Body:     "",
		}

		comment, err := service.CreateComment(ctx, cmd)
		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Contains(t, err.Error(), "body cannot be empty")
	})
}
