This is the roadmap from "Hello World" to a ship-ready MVP.

This plan follows the **"Steel Thread" methodology**. Instead of building "All the Backend" then "All the Frontend," we will build one complete feature at a time (Database  API  UI) so you always have a working product.

### Phase 1: The Foundation (Schema & Auth)

*Goal: A user can log in and see their profile. This sets up the critical "User Identity" context needed for everything else.*

| Step | Feature | Description | Key Tech |
| --- | --- | --- | --- |
| **1.1** | **Database Schema (Users)** | Create migration `001_users.sql`. Table: `users` (id, email, role, created_at). | SQL, Postgres |
| **1.2** | **OIDC Backend Logic** | Implement `GET /auth/login` (redirect to Google) and `GET /auth/callback` (exchange code for token). Create/Update user in DB on login. | Go, `coreos/go-oidc` |
| **1.3** | **Session Management** | Store the User ID in a secure HTTP-only cookie after login. Add middleware `RequireAuth` to protect API routes. | Go, `alexedwards/scs` (or Redis) |
| **1.4** | **Frontend Context** | Create a React Context (`AuthProvider`) that fetches `GET /api/me` on load. Show "Login with Google" button if 401. | React, Context API |
| **1.5** | **RBAC Middleware** | Implement `RequireRole("admin")` middleware. Verify you can lock down a dummy route. | Go Middleware |

### Phase 2: The Core Ticket Loop (Staff View)

*Goal: Staff can create, view, and close tickets. We ignore the "Public" aspect for a moment to get the data model right.*

| Step | Feature | Description | Key Tech |
| --- | --- | --- | --- |
| **2.1** | **Database Schema (Tickets)** | Migration `002_tickets.sql`. Tables: `tickets` (title, desc, status, priority, location), `ticket_events` (audit log). | SQL |
| **2.2** | **Ticket CRUD API** | Implement `POST /tickets`, `GET /tickets`, `GET /tickets/{id}`, `PATCH /tickets/{id}`. | Go, Chi, pgx |
| **2.3** | **Staff Dashboard UI** | Create a data table (TanStack Table recommended) to list tickets. Add sorting/filtering. | React, TanStack Table |
| **2.4** | **Ticket Detail View** | A page showing full ticket details. Allow changing Status (Open  Closed) and Priority. | React |
| **2.5** | **Audit Logging** | Hook into the `UpdateTicket` service. Every time a ticket changes, write a row to `ticket_events` (Who, When, What changed). | Go (Service Layer) |

### Phase 3: The Public Portal (The "Parishioner" View)

*Goal: Anonymous or authenticated public users can report issues without seeing sensitive data.*

| Step | Feature | Description | Key Tech |
| --- | --- | --- | --- |
| **3.1** | **Public API Layer** | Create a specific endpoint `GET /api/public/tickets` that *forces* filters (removes sensitive tickets, hides assignee names). | Go (Privacy Logic) |
| **3.2** | **Deduplication Search** | Build the "Search before you submit" UI. As user types "leak", query the public API and show matches. | React, Debounce |
| **3.3** | **Public Submission Form** | A simplified form for reporting. Includes "Is this sensitive?" checkbox. | React Hook Form |
| **3.4** | **"My Tickets" View** | If logged in as "Public", show a list of tickets *I* reported, even if they are marked sensitive/private. | SQL (`WHERE reporter_id = ?`) |

### Phase 4: Workflow & Assignments

*Goal: Moving the ticket from "New" to "Done" with the right people involved.*

| Step | Feature | Description | Key Tech |
| --- | --- | --- | --- |
| **4.1** | **Assignment Schema** | Add `assignee_user_id` (FK) and `assignee_name` (Text) to `tickets` table. | SQL |
| **4.2** | **Assignment UI** | On Ticket Detail, add a "Assign" dropdown. It should list Staff users AND allow typing a text string. | React, Combobox |
| **4.3** | **Comments System** | Migration `003_comments.sql`. API to add comments. Distinguish between "Internal Note" vs "Public Comment". | SQL, Go |
| **4.4** | **River Worker Setup** | Initialize the `River` worker pool in `main.go`. Create a basic `JobArgs` struct for notifications. | Go, River |
| **4.5** | **Email Notification** | Create a River worker `SendEmailWorker`. When ticket is assigned, enqueue a job. | Go, SMTP |

### Phase 5: MVP Polish & Data Sovereignty

*Goal: Production readiness.*

| Step | Feature | Description | Key Tech |
| --- | --- | --- | --- |
| **5.1** | **CSV Export** | Implement `GET /api/admin/export/tickets`. Stream the CSV response (don't load all into memory). | Go `encoding/csv` |
| **5.2** | **Status Configuration** | A simple Admin page to define valid statuses (e.g., "Waiting on Parts") and map them to Open/Resolved. | React, API |
| **5.3** | **Production Build** | Update `Dockerfile` to use a multi-stage build (Node build  Go build  Distroless image). | Docker |
| **5.4** | **Chaos Test** | Spin up the stack, kill the DB container, ensure the App reconnects. Kill the App, ensure no data loss in DB. | Manual QA |

### Recommended Order of Operations

**Start with Step 1.1 (User Schema).**
Do not try to build the Ticket form until you have a logged-in user to attach it to.

Would you like me to generate the **SQL for Step 1.1** (User Schema) to get you started?
