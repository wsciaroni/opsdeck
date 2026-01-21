ALTER TABLE organizations ADD COLUMN share_link_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE organizations ADD COLUMN share_link_token TEXT UNIQUE;
CREATE INDEX idx_organizations_share_link_token ON organizations(share_link_token);
