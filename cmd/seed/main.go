package main

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/storage/postgres"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	// 1. Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}
	seedEmail := os.Getenv("SEED_EMAIL")
	if seedEmail == "" {
		seedEmail = "seed@opsdeck.dev"
	}

	// 2. Connect to DB
	pool, err := postgres.ConnectPostgres(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// 3. Initialize Repositories
	userRepo := postgres.NewUserRepository(pool)
	orgRepo := postgres.NewOrganizationRepository(pool)
	ticketRepo := postgres.NewTicketRepository(pool)

	log.Printf("Starting seed process for user: %s", seedEmail)

	// 4. Find or Create User
	user, err := userRepo.GetByEmail(ctx, seedEmail)
	if err != nil {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	if user == nil {
		log.Printf("User not found. Creating user %s...", seedEmail)
		user = &domain.User{
			Email:     seedEmail,
			Name:      "Seed User",
			Role:      domain.RoleManager,
			AvatarURL: "",
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		log.Printf("User created with ID: %s", user.ID)
	} else {
		log.Printf("User found with ID: %s", user.ID)
	}

	// 5. Find or Create Organization
	memberships, err := orgRepo.ListByUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to list user memberships: %w", err)
	}

	var orgID uuid.UUID
	if len(memberships) > 0 {
		org := memberships[0].Organization
		orgID = org.ID
		log.Printf("Using existing organization: %s (%s)", org.Name, org.ID)
	} else {
		log.Printf("No organization found. Creating 'Seeder Corp'...")
		randomSuffix, err := randomHex(4)
		if err != nil {
			return fmt.Errorf("failed to generate random hex: %w", err)
		}
		slug := fmt.Sprintf("seeder-corp-%s", randomSuffix)
		org := &domain.Organization{
			Name: "Seeder Corp",
			Slug: slug,
		}
		if err := orgRepo.Create(ctx, org); err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}
		orgID = org.ID
		log.Printf("Organization created with ID: %s", orgID)

		// Add user as owner
		if err := orgRepo.AddMember(ctx, org.ID, user.ID, "owner"); err != nil {
			return fmt.Errorf("failed to add user to organization: %w", err)
		}
		log.Printf("User added as owner to organization")
	}

	// 6. Generate Tickets
	statuses := []string{
		domain.TicketStatusNew,
		domain.TicketStatusInProgress,
		domain.TicketStatusOnHold,
		domain.TicketStatusDone,
		domain.TicketStatusCanceled,
	}
	priorities := []string{
		domain.TicketPriorityLow,
		domain.TicketPriorityMedium,
		domain.TicketPriorityHigh,
		domain.TicketPriorityCritical,
	}
	locations := []string{
		"Server Room",
		"Main Office",
		"Warehouse",
		"Remote",
	}

	log.Printf("Generating 20 tickets for Organization %s...", orgID)
	for i := 1; i <= 20; i++ {
		status := statuses[rand.Intn(len(statuses))]
		priority := priorities[rand.Intn(len(priorities))]
		location := locations[rand.Intn(len(locations))]

		var assigneeID *uuid.UUID
		if rand.Float64() < 0.5 {
			// 50% chance to be assigned to seed user
			id := user.ID
			assigneeID = &id
		}

		var completedAt *time.Time
		if status == domain.TicketStatusDone || status == domain.TicketStatusCanceled {
			now := time.Now()
			completedAt = &now
		}

		ticket := &domain.Ticket{
			OrganizationID: orgID,
			Title:          fmt.Sprintf("Fix issue #%d with server", i),
			Description:    "Auto-generated description for testing purposes.",
			Location:       location,
			StatusID:       status,
			PriorityID:     priority,
			ReporterID:     user.ID,
			AssigneeUserID: assigneeID,
			CompletedAt:    completedAt,
		}

		if err := ticketRepo.Create(ctx, ticket); err != nil {
			return fmt.Errorf("failed to create ticket %d: %w", i, err)
		}
	}

	log.Printf("Seeded 20 tickets for Organization %s", orgID)
	return nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
