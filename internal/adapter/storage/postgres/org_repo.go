package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (name, slug)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, org.Name, org.Slug).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error {
	query := `
		INSERT INTO organization_members (organization_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, orgID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to add member to organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserMembership, error) {
	query := `
		SELECT o.id, o.name, o.slug, o.created_at, o.updated_at, om.role
		FROM organizations o
		JOIN organization_members om ON o.id = om.organization_id
		WHERE om.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations for user: %w", err)
	}
	defer rows.Close()

	var memberships []domain.UserMembership
	for rows.Next() {
		var m domain.UserMembership
		err := rows.Scan(
			&m.Organization.ID,
			&m.Organization.Name,
			&m.Organization.Slug,
			&m.Organization.CreatedAt,
			&m.Organization.UpdatedAt,
			&m.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user membership: %w", err)
		}
		memberships = append(memberships, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return memberships, nil
}

func (r *OrganizationRepository) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.Member, error) {
	query := `
		SELECT u.id, u.email, u.name, u.avatar_url, om.role
		FROM users u
		JOIN organization_members om ON u.id = om.user_id
		WHERE om.organization_id = $1
	`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization members: %w", err)
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		var avatarURL *string
		err := rows.Scan(
			&m.UserID,
			&m.Email,
			&m.Name,
			&avatarURL,
			&m.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		if avatarURL != nil {
			m.AvatarURL = *avatarURL
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return members, nil
}

func (r *OrganizationRepository) RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM organization_members
		WHERE organization_id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(ctx, query, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member from organization: %w", err)
	}
	return nil
}
