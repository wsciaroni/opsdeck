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

	// Run migrations
	if err := postgres.RunMigrations(ctx, dbURL); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skip("database not available")
	}
	defer pool.Close()

	// Clear tables
	_, err = pool.Exec(ctx, "TRUNCATE tickets, users, organizations CASCADE")
	require.NoError(t, err)

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

	t.Run("List By OrganizationIDs", func(t *testing.T) {
		filter := port.TicketFilter{
			OrganizationIDs: []uuid.UUID{org1.ID},
		}
		tickets, err := repo.List(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, tickets, 1)
		assert.Equal(t, ticket1.ID, tickets[0].ID)
	})

	t.Run("Sort By Priority DESC", func(t *testing.T) {
		tCritical := &domain.Ticket{OrganizationID: org1.ID, ReporterID: user.ID, Title: "Critical", StatusID: "new", PriorityID: "critical"}
		tHigh := &domain.Ticket{OrganizationID: org1.ID, ReporterID: user.ID, Title: "High", StatusID: "new", PriorityID: "high"}
		tMedium := &domain.Ticket{OrganizationID: org1.ID, ReporterID: user.ID, Title: "Medium", StatusID: "new", PriorityID: "medium"}

		require.NoError(t, repo.Create(ctx, tCritical))
		require.NoError(t, repo.Create(ctx, tHigh))
		require.NoError(t, repo.Create(ctx, tMedium))

		filter := port.TicketFilter{
			OrganizationID: &org1.ID,
			SortBy:         "priority",
			SortOrder:      "desc",
			Keyword:        nil, // Ensure no keyword search
		}
		tickets, err := repo.List(ctx, filter)
		require.NoError(t, err)

		// Filter out the initial ticket1 (Low) to make assertions easier, or include it.
		// Let's just check the first 3 are Critical, High, Medium.
		// Note: ticket1 is Low, so it should be last.

		// We expect 4 tickets in total for org1.
		require.Len(t, tickets, 4)

		assert.Equal(t, "critical", tickets[0].PriorityID)
		assert.Equal(t, "high", tickets[1].PriorityID)
		assert.Equal(t, "medium", tickets[2].PriorityID)
		assert.Equal(t, "low", tickets[3].PriorityID)
	})

	t.Run("Sort By CreatedAt ASC", func(t *testing.T) {
		filter := port.TicketFilter{
			OrganizationID: &org1.ID,
			SortBy:         "created_at",
			SortOrder:      "asc",
		}
		tickets, err := repo.List(ctx, filter)
		require.NoError(t, err)
		require.Len(t, tickets, 4)

		// ticket1 was created first (in setup).
		assert.Equal(t, ticket1.ID, tickets[0].ID)
	})
}
