-- Add index for dashboard ticket list query
-- Speeds up: WHERE organization_id = ? ORDER BY created_at DESC
CREATE INDEX idx_tickets_org_created_at ON tickets(organization_id, created_at DESC);
