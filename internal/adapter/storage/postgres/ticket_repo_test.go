package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wsciaroni/opsdeck/internal/adapter/storage/postgres"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

func TestTicketRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skip("database not available")
	}
	defer pool.Close()

	// Clear tables
	pool.Exec(ctx, "TRUNCATE tickets, users, organizations CASCADE")

	repo := postgres.NewTicketRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	orgRepo := postgres.NewOrganizationRepository(pool)

	// Setup Data
	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Role:      domain.RoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	org1 := &domain.Organization{
		ID:        uuid.New(),
		Name:      "Org 1",
		Slug:      "org-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, orgRepo.Create(ctx, org1))

	org2 := &domain.Organization{
		ID:        uuid.New(),
		Name:      "Org 2",
		Slug:      "org-2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, orgRepo.Create(ctx, org2))

	ticket1 := &domain.Ticket{
		OrganizationID: org1.ID,
		ReporterID:     user.ID,
		Title:          "Ticket Org 1",
		Description:    "Desc 1",
		StatusID:       "new",
		PriorityID:     "low",
		Location:       "Loc 1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, repo.Create(ctx, ticket1))

	ticket2 := &domain.Ticket{
		OrganizationID: org2.ID,
		ReporterID:     user.ID,
		Title:          "Ticket Org 2",
		Description:    "Desc 2",
		StatusID:       "new",
		PriorityID:     "low",
		Location:       "Loc 2",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, repo.Create(ctx, ticket2))

	t.Run("List By Organization", func(t *testing.T) {
		filter := port.TicketFilter{
			OrganizationID: &org1.ID,
		}
		tickets, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, tickets, 1)
		assert.Equal(t, ticket1.ID, tickets[0].ID)
	})

	t.Run("List All (Admin Export)", func(t *testing.T) {
		filter := port.TicketFilter{
			OrganizationID: nil,
		}
		tickets, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, tickets, 2)
	})
}
