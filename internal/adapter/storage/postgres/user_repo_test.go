package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

func TestUserRepository(t *testing.T) {
	// Skip if no DB available (optional, but good practice if not always running with DB)
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := ConnectPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	// Run migration manually
	// Read file
	migrationSQL, err := os.ReadFile("../../../../migrations/001_users.sql")
	if err != nil {
		t.Fatalf("Failed to read migration file: %v", err)
	}
	_, err = pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		t.Fatalf("Failed to run migration: %v", err)
	}

	// Cleanup
	defer func() {
		pool.Exec(ctx, "TRUNCATE TABLE users")
	}()

	repo := NewUserRepository(pool)

	// Test Create
	user := &domain.User{
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      domain.RoleStaff,
		AvatarURL: "http://example.com/avatar.png",
	}

	err = repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Errorf("Expected ID to be set, got Nil")
	}
	if user.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be set")
	}

	// Test GetByID
	fetchedUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetchedUser == nil {
		t.Fatalf("Expected user found, got nil")
	}
	if fetchedUser.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, fetchedUser.ID)
	}
	if fetchedUser.Email != user.Email {
		t.Errorf("Expected Email %v, got %v", user.Email, fetchedUser.Email)
	}

	// Test GetByEmail
	fetchedUserEmail, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if fetchedUserEmail == nil {
		t.Fatalf("Expected user found by email, got nil")
	}
	if fetchedUserEmail.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, fetchedUserEmail.ID)
	}

	// Test Update
	user.Name = "Updated Name"
	user.Role = domain.RoleAdmin
	time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt changes
	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updatedUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updatedUser.Name != "Updated Name" {
		t.Errorf("Expected Name 'Updated Name', got '%v'", updatedUser.Name)
	}
	if updatedUser.Role != domain.RoleAdmin {
		t.Errorf("Expected Role 'admin', got '%v'", updatedUser.Role)
	}
	if !updatedUser.UpdatedAt.After(user.CreatedAt) {
		t.Errorf("Expected UpdatedAt to be after CreatedAt")
	}

	// Test NotFound
	missingUser, err := repo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID missing failed (should not error): %v", err)
	}
	if missingUser != nil {
		t.Errorf("Expected nil for missing user, got %v", missingUser)
	}

	missingUserEmail, err := repo.GetByEmail(ctx, "missing@example.com")
	if err != nil {
		t.Fatalf("GetByEmail missing failed (should not error): %v", err)
	}
	if missingUserEmail != nil {
		t.Errorf("Expected nil for missing email, got %v", missingUserEmail)
	}
}
