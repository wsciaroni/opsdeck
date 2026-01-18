# OpsDeck Use Cases

**Scope:** This document outlines the functional requirements for OpsDeck.
**MVP Status:**

* ✅ = **MVP (Phase 1):** Essential for launch.
* ❌ = **Post-MVP (Phase 2+):** Enhancements for scale, automation, and advanced asset management.

## 1. Authentication & User Management

*Core security and role definition.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-01** | ✅ | **Login via SSO** | All Users | Authenticate using a third-party provider (Google/OIDC). |
| **UC-02** | ✅ | **Auto-Provision Account** | System | Automatically create a "Public" level account upon first SSO login. |
| **UC-03** | ✅ | **Assign User Roles** | Admin | Promote a user from "Public" to "Staff" or "Manager." |
| **UC-04** | ✅ | **Revoke Access** | Admin | Disable a user’s access to the system. |

## 2. Public Portal (Reporting)

*The interface for parishioners and members.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-05** | ✅ | **Search Public Tickets** | Public | Search existing tickets by keyword to prevent duplicates. |
| **UC-06** | ✅ | **View Public Ticket Details** | Public | View status and public comments of non-sensitive tickets. |
| **UC-07** | ✅ | **Submit New Ticket** | Public | Fill out a form (Location, Category, Description, Photo) to report an issue. |
| **UC-08** | ✅ | **Flag as Sensitive** | Public | User checks "Sensitive" during submission to hide the ticket from the public feed. |
| **UC-09** | ✅ | **Receive Status Updates** | Public | Receive email notifications when the status of a reported ticket changes. |
| **UC-75** | ✅ | **Submit on Behalf Of** | Staff | Staff submits a ticket but manually sets a different user as the "Reporter" (e.g., phone call intake). |

## 3. Ticket Management (Work)

*The core workflow for staff and managers.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-10** | ✅ | **View All Tickets** | Staff | View list of *all* tickets (including sensitive) with filtering options. |
| **UC-11** | ✅ | **Update Ticket Status** | Staff | Change status (e.g., New  In Progress  Done). |
| **UC-12** | ✅ | **Change Visibility (Hide)** | Manager | Change a ticket from "Public" to "Sensitive" (and vice versa). |
| **UC-14** | ✅ | **Add Internal Note** | Staff | Add a comment visible *only* to Staff/Managers. |
| **UC-15** | ✅ | **Add Public Comment** | Staff | Add a comment visible to the reporter and public. |
| **UC-16** | ✅ | **Assign Ticket (User)** | Manager | Assign a ticket to a registered Staff member. |
| **UC-32** | ✅ | **Assign Ticket (Text)** | Manager | Assign a ticket to a non-user string (e.g., "Bob's Plumbing"). |
| **UC-33** | ✅ | **Hide Assignment Details** | System | Ensure "Assigned To" field is never shown to the public. |
| **UC-66** | ❌ | **Merge Duplicates** | Manager | Combine multiple tickets into one parent ticket. |
| **UC-67** | ❌ | **Reopen Ticket** | Manager | Change status from Resolved back to Open. |

## 4. Workflow Configuration

*Defining how the work flows.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-24** | ✅ | **Manage Custom Statuses** | Admin | Create statuses (e.g., "Waiting on Budget") and map them to System States (Open/Closed). |
| **UC-25** | ✅ | **Manage Categories** | Admin | Define locations (Hall, Kitchen) or Types (Plumbing). |
| **UC-55** | ✅ | **Define Priority Levels** | Admin | Create priority levels (Low, Med, Critical) with custom descriptions. |
| **UC-37** | ✅ | **Filter Resolved Publicly** | System | Automatically hide tickets mapped to "Resolved" from public search. |

## 5. Planning & Organization

*Tools for managing the backlog.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-17** | ✅ | **Filter & Sort Backlog** | Staff | Filter by Location, Priority, or Status. |
| **UC-18** | ❌ | **Kanban Board View** | Staff | View tickets in drag-and-drop columns. |
| **UC-19** | ❌ | **Bulk Edit Tickets** | Manager | Select multiple tickets to close or move them. |
| **UC-73** | ❌ | **Calendar Subscription** | Staff | iCal feed for assigned tickets. |

## 6. Teams & Notifications

*Routing work to the right people.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-47** | ✅ | **Global Notification Prefs** | Admin | Set default email settings for the organization. |
| **UC-49** | ✅ | **User Notification Prefs** | User | Users override defaults for their own email/SMS. |
| **UC-40** | ❌ | **Create Teams** | Admin | Define groups (e.g., "Electrical Team"). |
| **UC-43** | ❌ | **Dynamic Routing Rules** | Admin | "If Category=Plumbing, Assign to Team A." |
| **UC-48** | ❌ | **Team Notification Prefs** | Team Lead | Set notification defaults for a specific team. |

## 7. Recurring Maintenance

*Scheduled tasks.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-21** | ❌ | **Create Recurring Schedule** | Manager | Define template and cadence (e.g., "Air Filters" every 3 months). |
| **UC-22** | ❌ | **Generate Recurring Ticket** | System | Auto-create ticket when due date arrives. |

## 8. Data Sovereignty (Export)

*Preventing vendor lock-in.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-26** | ✅ | **Export Ticket Data (CSV)** | Admin | Download spreadsheet of all tickets (past and present). |
| **UC-27** | ✅ | **Export User Directory** | Admin | Download list of all users and roles. |
| **UC-28** | ❌ | **Batch Download Attachments** | Admin | Download ZIP of all photos/docs. |

## 9. Assets & Inventory

*Tracking physical equipment.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-78** | ❌ | **Register Asset** | Admin | Create record for equipment (e.g., HVAC #1) with serial number. |
| **UC-79** | ❌ | **Link Ticket to Asset** | Staff | Associate a repair with a specific asset history. |
| **UC-81** | ❌ | **Print QR Codes** | Admin | Generate QR labels for assets to quick-start tickets. |

## 10. Financials & Audit

*Accountability and costs.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-63** | ❌ | **View Audit Log** | Manager | Detailed history of who changed what field. |
| **UC-70** | ❌ | **Log Labor/Material Costs** | Staff | Record time and money spent on a ticket. |
| **UC-72** | ❌ | **Generate Cost Report** | Manager | Report on spending by category. |

## 11. System Maintenance

*DevOps and Admin.*

| ID | MVP | Use Case Name | Actors | Description |
| --- | --- | --- | --- | --- |
| **UC-88** | ❌ | **Trigger Manual Backup** | Admin | Force database snapshot. |
| **UC-89** | ✅ | **System Health Check** | Admin | Basic health endpoint for monitoring (e.g., `/healthz`). |
