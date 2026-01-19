-- Create Ticket Statuses Table
CREATE TABLE ticket_statuses (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    level INT NOT NULL
);

INSERT INTO ticket_statuses (id, label, level) VALUES
('new', 'New', 1),
('in_progress', 'In Progress', 2),
('on_hold', 'On Hold', 3),
('done', 'Done', 4),
('canceled', 'Canceled', 5);

-- Create Ticket Priorities Table
CREATE TABLE ticket_priorities (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    level INT NOT NULL
);

INSERT INTO ticket_priorities (id, label, level) VALUES
('low', 'Low', 1),
('medium', 'Medium', 2),
('high', 'High', 3),
('critical', 'Critical', 4);

-- Create Tickets Table
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    reporter_id UUID NOT NULL REFERENCES users(id),
    assignee_user_id UUID REFERENCES users(id),
    status_id TEXT NOT NULL REFERENCES ticket_statuses(id),
    priority_id TEXT NOT NULL REFERENCES ticket_priorities(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_tickets_org_status ON tickets(organization_id, status_id);
CREATE INDEX idx_tickets_reporter ON tickets(reporter_id);
CREATE INDEX idx_tickets_assignee ON tickets(assignee_user_id);
