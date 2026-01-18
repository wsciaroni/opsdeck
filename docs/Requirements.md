### 1. Functional Requirements (FR)

These define *what* the system must do.

#### 1.1. Authentication & Roles

* **FR-01 (SSO):** The system shall support authentication via OpenID Connect (OIDC), specifically allowing Google Sign-In.
* **FR-02 (Auto-Provisioning):** Upon first successful login, the system shall automatically create a user account with the default role of "Public."
* **FR-03 (RBAC):** The system shall enforce Role-Based Access Control with the following hierarchy:
* **Public:** Can search public tickets, submit new tickets, and view tickets they reported.
* **Staff:** Can view all tickets (including sensitive), assign tickets, change status, and add internal notes.
* **Manager/Admin:** Can manage users, configure system settings (categories, priorities), and export data.



#### 1.2. Ticket Lifecycle

* **FR-04 (Submission):** The system shall accept ticket submissions requiring: Category, Priority, Location, Title, and Description. Photo upload must be optional but supported.
* **FR-05 (Deduplication Search):** The ticket submission page shall perform a live search of existing *Open* and *Public* tickets to suggest duplicates to the user before they submit.
* **FR-06 (Visibility Toggling):** The system shall allow a Manager to toggle a ticket's visibility between "Public" and "Sensitive" at any time.
* **FR-07 (Status Mapping):** The system shall allow Admins to define custom statuses (e.g., "Waiting on Parts") and map them to one of two system states: `Open` (Active) or `Resolved` (Closed).

#### 1.3. Workflow & Assignment

* **FR-08 (Dual Assignment):** The system shall allow a ticket to be assigned to either:
* A registered **User** (Staff member).
* A **Plain Text String** (e.g., "Bob's Roofing") for external vendors.


* **FR-09 (Privacy Enforcement):** The system shall strictly redact the "Assigned To," "Internal Notes," and "Vendor Info" fields from all API responses sent to Public users.

#### 1.4. Notifications

* **FR-10 (Triggers):** The system shall generate notifications for the following events: `TicketCreated`, `TicketAssigned`, `TicketCommented`, `StatusChanged`.
* **FR-11 (Preference Hierarchy):** The system shall determine notification delivery (Email vs. SMS vs. None) by checking preferences in this specific order:
1. User Preference (if set)
2. Team Preference (if set)
3. Global Default



#### 1.5. Data Sovereignty (Export)

* **FR-12 (CSV Export):** The system shall provide an admin endpoint to generate and download a `.csv` file containing the full history of tickets.
* **FR-13 (User Export):** The system shall provide an admin endpoint to export the user directory to `.csv` to facilitate migration.

---

### 2. Non-Functional Requirements (NFR)

These define *how* the system performs and behaves.

#### 2.1. Architecture & Deployment

* **NFR-01 (Statelessness):** The application backend must be stateless to survive container restarts (Chaos Monkey). Session data must be stored in an external store (Redis) or use stateless tokens (JWT).
* **NFR-02 (Single Artifact):** The build process must result in a single deployable artifact (Docker container or Binary) that contains both the Go backend and the compiled React frontend.
* **NFR-03 (Database Portability):** The system must use PostgreSQL (v16+) as the primary data store and manage schema changes via automated migration scripts on startup.

#### 2.2. Performance & Reliability

* **NFR-04 (Search Latency):** Public ticket search queries must return results in under 200ms for datasets up to 100,000 tickets.
* **NFR-05 (Job Reliability):** Background tasks (e.g., sending emails) must be processed via a persistent queue (Postgres-based or Redis-based) with automatic retry logic for failed jobs.

#### 2.3. Security

* **NFR-06 (Data Minimization):** The system shall not expose user email addresses or phone numbers in any public API response.
* **NFR-07 (Audit Trail):** All state-changing actions (Status Change, Assignment, Configuration Edit) must be logged to an immutable audit table with `UserID`, `Timestamp`, `Action`, and `PreviousValue`.

---

### 3. Interface Requirements

* **UI-01 (Responsiveness):** The web interface must be fully functional on mobile devices (320px width minimum) to support staff working in the field.
* **UI-02 (Accessibility):** The public-facing reporting portal must comply with WCAG 2.1 AA standards (color contrast, keyboard navigation) to ensure accessibility for all parishioners.
