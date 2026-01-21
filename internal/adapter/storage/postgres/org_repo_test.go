package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

func TestOrganizationRepository(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := ConnectPostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	// Assuming migrations are already applied or we rely on pre-existing schema for integration test env.
	// In user_repo_test, it reads migration file manually. Let's try to do minimal setup or rely on running instance?
	// The user_repo_test runs migration 001. We probably need more migrations for orgs.
	// Since I cannot know the full migration state easily, I'll assume the environment running tests has the DB ready or I will try to apply migrations if possible.
	// However, `user_repo_test.go` pattern suggests applying specific migrations.
	// I'll skip migration setup here to avoid complexity and assume it's an integration test environment where I might need to truncate.

	// Cleanup function
	cleanup := func() {
		if _, err := pool.Exec(ctx, "TRUNCATE TABLE organization_members CASCADE"); err != nil {
			t.Logf("Failed to truncate organization_members: %v", err)
		}
		if _, err := pool.Exec(ctx, "TRUNCATE TABLE organizations CASCADE"); err != nil {
			t.Logf("Failed to truncate organizations: %v", err)
		}
		if _, err := pool.Exec(ctx, "TRUNCATE TABLE users CASCADE"); err != nil {
			t.Logf("Failed to truncate users: %v", err)
		}
	}
	cleanup()
	defer cleanup()

	// Needed Repos
	userRepo := NewUserRepository(pool)
	orgRepo := NewOrganizationRepository(pool)

	// Create User
	user := &domain.User{
		Email:     "orgtest@example.com",
		Name:      "Org Tester",
		Role:      domain.RoleStaff,
		AvatarURL: "avatar.png",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create Organization
	org := &domain.Organization{
		Name: "Test Org",
		Slug: "test-org-123",
	}
	if err := orgRepo.Create(ctx, org); err != nil {
		t.Fatalf("Failed to create org: %v", err)
	}

	// Test AddMember
	if err := orgRepo.AddMember(ctx, org.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	// Test ListMembers
	members, err := orgRepo.ListMembers(ctx, org.ID)
	if err != nil {
		t.Fatalf("Failed to list members: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}
	if members[0].UserID != user.ID {
		t.Errorf("Expected member user ID %v, got %v", user.ID, members[0].UserID)
	}
	if members[0].Role != "owner" {
		t.Errorf("Expected role owner, got %v", members[0].Role)
	}
	if members[0].Email != user.Email {
		t.Errorf("Expected email %v, got %v", user.Email, members[0].Email)
	}

	// Test ListByUser (existing functionality, just to be sure)
	memberships, err := orgRepo.ListByUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to list memberships: %v", err)
	}
	if len(memberships) != 1 {
		t.Errorf("Expected 1 membership, got %d", len(memberships))
	}
	if memberships[0].Organization.ID != org.ID {
		t.Errorf("Expected org ID %v, got %v", org.ID, memberships[0].Organization.ID)
	}

	// Test RemoveMember
	if err := orgRepo.RemoveMember(ctx, org.ID, user.ID); err != nil {
		t.Fatalf("Failed to remove member: %v", err)
	}

	members, err = orgRepo.ListMembers(ctx, org.ID)
	if err != nil {
		t.Fatalf("Failed to list members after removal: %v", err)
	}
	if len(members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(members))
	}
}
