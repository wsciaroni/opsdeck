This is the complete **Phase 1 & 2 Database Schema** for OpsDeck.

It is designed for PostgreSQL 16+ and uses `UUIDs` for primary keys (essential for data merging and distributed systems). It also leverages `JSONB` for the audit logs to keep them flexible.

### File: `migrations/001_initial_schema.sql`

```sql
-- Enable UUID extension for ID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. Configuration Tables (Lookups)
-- ==========================================

-- Priorities: Defines urgency levels (Low, Medium, Critical)
CREATE TABLE priorities (
    id VARCHAR(50) PRIMARY KEY, -- e.g., 'critical', 'low'
    label VARCHAR(100) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL, -- For sorting (10=Low, 90=Critical)
    is_default BOOLEAN DEFAULT FALSE
);

-- Statuses: Defines lifecycle states (Open, Waiting, Done)
CREATE TABLE statuses (
    id VARCHAR(50) PRIMARY KEY, -- e.g., 'in_progress', 'waiting_parts'
    label VARCHAR(100) NOT NULL,
    category VARCHAR(20) NOT NULL CHECK (category IN ('open', 'resolved')), -- System logic mapping
    is_default BOOLEAN DEFAULT FALSE
);

-- Seed Default Data
INSERT INTO priorities (id, label, description, level, is_default) VALUES
('low', 'Low', 'Cosmetic issues or non-urgent repairs.', 10, FALSE),
('medium', 'Medium', 'Standard maintenance affecting daily use.', 50, TRUE),
('high', 'High', 'Urgent issues affecting safety or operations.', 80, FALSE),
('critical', 'Critical', 'Emergency: Fire, Flood, or Security Risk.', 100, FALSE);

INSERT INTO statuses (id, label, category, is_default) VALUES
('new', 'New', 'open', TRUE),
('in_progress', 'In Progress', 'open', FALSE),
('waiting', 'Waiting', 'open', FALSE),
('resolved', 'Resolved', 'resolved', FALSE),
('wont_fix', 'Won''t Fix', 'resolved', FALSE);


-- ==========================================
-- 2. User Management
-- ==========================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('public', 'staff', 'manager', 'admin')),
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 3. Core Ticket Data
-- ==========================================

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Ownership
    reporter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Assignment (Dual Mode: specific user OR text string)
    assignee_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assignee_name VARCHAR(255), -- "Bob's Plumbing" or fallback if user is null
    
    -- Content
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    
    -- State
    status_id VARCHAR(50) NOT NULL REFERENCES statuses(id),
    priority_id VARCHAR(50) NOT NULL REFERENCES priorities(id),
    
    -- Visibility
    is_public BOOLEAN NOT NULL DEFAULT TRUE, -- False = Sensitive/Hidden
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Indexes for common lookups
CREATE INDEX idx_tickets_status ON tickets(status_id);
CREATE INDEX idx_tickets_public_search ON tickets(is_public) WHERE is_public = TRUE;
CREATE INDEX idx_tickets_reporter ON tickets(reporter_id);
CREATE INDEX idx_tickets_assignee ON tickets(assignee_user_id);


-- ==========================================
-- 4. Communication & Audit
-- ==========================================

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    is_internal BOOLEAN NOT NULL DEFAULT FALSE, -- True = Staff Only
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ticket_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    event_type VARCHAR(50) NOT NULL, -- 'created', 'status_change', 'assigned', 'comment'
    
    -- JSONB allows us to store "Old vs New" flexibly without strict schema
    -- Example: {"field": "status", "old": "new", "new": "in_progress"}
    changes JSONB NOT NULL DEFAULT '{}'::jsonb, 
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 5. Notifications
-- ==========================================

CREATE TABLE notification_preferences (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- NULL = Global Default
    event_type VARCHAR(50) NOT NULL, -- 'ticket_created', 'ticket_assigned'
    
    email_enabled BOOLEAN DEFAULT TRUE,
    sms_enabled BOOLEAN DEFAULT FALSE,
    
    -- Enforce one row per user per event type
    -- Use a partial unique index so we can have one NULL user_id (Global) per event
    UNIQUE (user_id, event_type)
);

-- ==========================================
-- 6. River Job Queue (Standard Schema)
-- ==========================================
-- Note: River requires its own migration, but this is the standard v5 schema
-- for reference. You will likely run this via the river CLI or library.

CREATE TABLE river_job (
    id bigserial PRIMARY KEY,
    state varchar(8) NOT NULL DEFAULT 'available',
    attempt smallint NOT NULL DEFAULT 0,
    max_attempts smallint NOT NULL,
    attempted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    finalized_at timestamptz,
    scheduled_at timestamptz NOT NULL DEFAULT NOW(),
    priority smallint NOT NULL DEFAULT 1,
    args jsonb,
    attempt_errors jsonb[],
    tags varchar(255)[],
    meta jsonb
);

-- River Indexes
CREATE INDEX river_job_state_priority_scheduled_at_id_idx ON river_job (state, priority, scheduled_at, id);

```

### Visualizing the Relationships

### Key Design Decisions

1. **`assignee_name` vs `assignee_user_id`:**
* This supports **UC-32** (Assign to Text).
* Logic: If `assignee_user_id` is present, display that User's Name. If it is null, display the `assignee_name` string.


2. **`statuses` Table (Not Enum):**
* We used a table instead of a Postgres ENUM.
* **Why:** This allows the Admin to go into the UI and create a new status called "Waiting on Contractor" without you (the developer) having to write a database migration.


3. **`ticket_events` (JSONB):**
* We use a JSONB column for the changes.
* **Why:** Audit logs are messy. Sometimes you change a Status (simple string), sometimes you change a Description (long text), sometimes you change an Assignment (UUID). JSONB handles this polymorphism perfectly without 50 empty columns.

## Example Go Models

Here are the **Go Domain Models** (`structs`) that map to your database schema.

In a **Modular Monolith**, these belong in **`internal/core/domain/`**. They serve as the common language between your Database Layer (Storage Adapter) and your API Layer (Web Adapter).

I have used pointers (e.g., `*uuid.UUID`) for nullable fields to ensure they serialize correctly to JSON as `null` rather than empty strings/zeros.

### 1. User Domain

**File:** `internal/core/domain/user.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RolePublic  Role = "public"
	RoleStaff   Role = "staff"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Role      Role      `json:"role" db:"role"`
	AvatarURL string    `json:"avatar_url,omitempty" db:"avatar_url"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

```

### 2. Configuration (Lookups)

**File:** `internal/core/domain/config.go`

```go
package domain

type StatusCategory string

const (
	StatusOpen     StatusCategory = "open"
	StatusResolved StatusCategory = "resolved"
)

type Priority struct {
	ID          string `json:"id" db:"id"`
	Label       string `json:"label" db:"label"`
	Description string `json:"description" db:"description"`
	Level       int    `json:"level" db:"level"` // 10=Low, 100=Critical
	IsDefault   bool   `json:"is_default" db:"is_default"`
}

type Status struct {
	ID        string         `json:"id" db:"id"`
	Label     string         `json:"label" db:"label"`
	Category  StatusCategory `json:"category" db:"category"` // "open" or "resolved"
	IsDefault bool           `json:"is_default" db:"is_default"`
}

```

### 3. Ticket Core

**File:** `internal/core/domain/ticket.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Location    string    `json:"location" db:"location"`

	// Relationships
	ReporterID uuid.UUID `json:"reporter_id" db:"reporter_id"`
	
	// Assignment (Dual Mode)
	// Pointer because it can be unassigned (null)
	AssigneeUserID *uuid.UUID `json:"assignee_user_id,omitempty" db:"assignee_user_id"` 
	AssigneeName   string     `json:"assignee_name,omitempty" db:"assignee_name"`

	// State
	StatusID   string `json:"status_id" db:"status_id"`
	PriorityID string `json:"priority_id" db:"priority_id"`

	// Visibility
	IsPublic bool `json:"is_public" db:"is_public"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// TicketDetail is a "Join View" useful for the API response 
// so the frontend doesn't have to fetch Users/Statuses separately.
type TicketDetail struct {
	Ticket
	ReporterName   string `json:"reporter_name" db:"reporter_name"`
	AssigneeAvatar string `json:"assignee_avatar,omitempty" db:"assignee_avatar"`
	StatusLabel    string `json:"status_label" db:"status_label"`
	PriorityLabel  string `json:"priority_label" db:"priority_label"`
}

```

### 4. Events & Comments (Audit Log)

**File:** `internal/core/domain/activity.go`

```go
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventCreated      EventType = "created"
	EventStatusChange EventType = "status_change"
	EventAssigned     EventType = "assigned"
	EventComment      EventType = "comment"
)

// Comment represents a text discussion on a ticket
type Comment struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	TicketID   uuid.UUID  `json:"ticket_id" db:"ticket_id"`
	AuthorID   *uuid.UUID `json:"author_id" db:"author_id"` // Null if system/deleted user
	Content    string     `json:"content" db:"content"`
	IsInternal bool       `json:"is_internal" db:"is_internal"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// TicketEvent represents the immutable audit log
type TicketEvent struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TicketID  uuid.UUID  `json:"ticket_id" db:"ticket_id"`
	ActorID   *uuid.UUID `json:"actor_id" db:"actor_id"`
	EventType EventType  `json:"event_type" db:"event_type"`
	
	// Changes is stored as JSONB in Postgres.
	// Map allows flexible structure: {"old": "open", "new": "closed"}
	Changes json.RawMessage `json:"changes" db:"changes"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

```

### Usage Note: `json.RawMessage`

For the `TicketEvent.Changes` field, I used `json.RawMessage`.

* **Why:** When reading from the DB, we don't want to force a strict schema on the "Changes" column because it varies wildly (Status changes look different than Description edits).
* **Effect:** The API will pass the raw JSON straight through to the Frontend, which can decide how to render it.
