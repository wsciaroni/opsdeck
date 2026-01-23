ALTER TABLE organizations ADD COLUMN public_view_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE organizations ADD COLUMN public_view_token TEXT UNIQUE;
CREATE INDEX idx_organizations_public_view_token ON organizations(public_view_token);
