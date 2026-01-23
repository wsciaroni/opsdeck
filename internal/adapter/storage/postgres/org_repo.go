package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, share_link_enabled, share_link_token, public_view_enabled, public_view_token, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	var org domain.Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.ShareLinkEnabled,
		&org.ShareLinkToken,
		&org.PublicViewEnabled,
		&org.PublicViewToken,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *OrganizationRepository) GetByShareToken(ctx context.Context, token string) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, share_link_enabled, share_link_token, public_view_enabled, public_view_token, created_at, updated_at
		FROM organizations
		WHERE share_link_token = $1
	`
	var org domain.Organization
	err := r.db.QueryRow(ctx, query, token).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.ShareLinkEnabled,
		&org.ShareLinkToken,
		&org.PublicViewEnabled,
		&org.PublicViewToken,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get organization by token: %w", err)
	}
	return &org, nil
}

func (r *OrganizationRepository) GetByPublicViewToken(ctx context.Context, token string) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, share_link_enabled, share_link_token, public_view_enabled, public_view_token, created_at, updated_at
		FROM organizations
		WHERE public_view_token = $1
	`
	var org domain.Organization
	err := r.db.QueryRow(ctx, query, token).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.ShareLinkEnabled,
		&org.ShareLinkToken,
		&org.PublicViewEnabled,
		&org.PublicViewToken,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get organization by public view token: %w", err)
	}
	return &org, nil
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (name, slug, share_link_enabled, share_link_token, public_view_enabled, public_view_token)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, org.Name, org.Slug, org.ShareLinkEnabled, org.ShareLinkToken, org.PublicViewEnabled, org.PublicViewToken).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, slug = $2, share_link_enabled = $3, share_link_token = $4, public_view_enabled = $5, public_view_token = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, org.Name, org.Slug, org.ShareLinkEnabled, org.ShareLinkToken, org.PublicViewEnabled, org.PublicViewToken, org.ID).Scan(&org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
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
		SELECT o.id, o.name, o.slug, o.share_link_enabled, o.share_link_token, o.public_view_enabled, o.public_view_token, o.created_at, o.updated_at, om.role
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
			&m.Organization.ShareLinkEnabled,
			&m.Organization.ShareLinkToken,
			&m.Organization.PublicViewEnabled,
			&m.Organization.PublicViewToken,
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

func (r *OrganizationRepository) UpdateMemberRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error {
	query := `
		UPDATE organization_members
		SET role = $1
		WHERE organization_id = $2 AND user_id = $3
	`
	_, err := r.db.Exec(ctx, query, role, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	return nil
}
