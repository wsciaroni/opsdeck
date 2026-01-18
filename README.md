# OpsDeck

<!--
[![Go Report Card](https://goreportcard.com/badge/github.com/wsciaroni/opsdeck)](https://goreportcard.com/report/github.com/wsciaroni/opsdeck)
[![Build Status](https://img.shields.io/github/actions/workflow/status/wsciaroni/opsdeck/ci.yml?branch=main)](https://github.com/wsciaroni/opsdeck/actions)
[![License](https://img.shields.io/github/license/wsciaroni/opsdeck)](LICENSE)
[![Docker Image](https://img.shields.io/docker/pulls/wsciaroni/opsdeck.svg)](https://hub.docker.com/r/wsciaroni/opsdeck)
-->

**Open-source facility maintenance tracker for organizations.**

Features public reporting with smart deduplication, granular RBAC, recurring tasks, and asset management. Built for data sovereignty with full export capabilities and strict privacy controls to prevent vendor lock-in.

---

## 📖 About

OpsDeck is a maintenance ticket tracking system designed to balance **public transparency** with **internal discretion**. It solves the problem of managing facilities with a mix of paid staff, volunteers, and public reporters, without the complexity or cost of enterprise ERP software.

### Core Philosophy
* **No Vendor Lock-in:** Comprehensive export of all data (Tickets, Users, Assets) to open formats (CSV/JSON) at any time.
* **Privacy First:** Granular control over what the public sees vs. what staff sees.
* **Chaos Resilient:** Stateless architecture designed to survive random instance failures.
* **Deployment Flexibility:** Runs on a single Docker container or a scalable Kubernetes cluster.

---

## 🚀 Key Features

* **Smart Public Portal:** Public users can search existing tickets to prevent duplicates before reporting.
* **Role-Based Access Control (RBAC):** Distinct views for Public, Staff, Team Owners, and Facility Managers.
* **Dynamic Routing:** Automatically route tickets to specific teams (e.g., "Electrical", "Plumbing") based on category or location rules.
* **Recurring Maintenance:** "Set and forget" schedules for routine tasks (e.g., HVAC filters, fire inspections).
* **Asset Management:** Track repair history against specific physical assets (QR code support).
* **Notification Cascade:** Configurable overrides for Global -> Team -> User notification preferences.
* **Audit Logging:** Complete history of who changed what and when.

---

## 🛠 Tech Stack

OpsDeck is built as a **Stateless Modular Monolith**. It compiles into a single binary for ease of deployment but supports horizontal scaling via external state management.

* **Backend:** Go (Golang) 1.24+ using Echo/Chi.
* **Frontend:** React + TypeScript (embedded into the Go binary).
* **Database:** PostgreSQL 16 (Relational data + JSONB).
* **Job Queue:** River (Go) running on PostgreSQL (Transactional reliability).
* **Resilience:** Redis (Optional, for distributed locking/sessions in scaled mode).

---

## 📦 Deployment

### Mode A: Simple (Single Container)
Ideal for small organizations. The frontend is embedded in the backend binary.

1.  **Clone and Start:**
    ```bash
    git clone [https://github.com/wsciaroni/opsdeck.git](https://github.com/wsciaroni/opsdeck.git)
    cd opsdeck
    docker compose up -d
    ```
2.  **Access:**
    * App: `http://localhost:8080`
    * Default Admin: `admin@example.com` / `password`

### Mode B: Scaled (Kubernetes/Cloud)
Designed to survive the "Chaos Monkey."

* **Stateless:** The app container contains no state. Sessions are stored in Redis; Jobs are stored in Postgres.
* **Horizontal Scaling:** Deploy 10+ replicas behind a Load Balancer.
* **Configuration:** See `k8s/` directory for Helm charts.

---

## 💻 Development

### Prerequisites
* Go 1.24+
* Node.js 20+
* PostgreSQL 16
* Docker (for local DB)

### Quick Start
1.  **Start Dependencies:**
    ```bash
    docker compose -f docker-compose.dev.yml up -d db
    ```
2.  **Run Backend (Hot Reload):**
    ```bash
    go run cmd/server/main.go
    ```
3.  **Run Frontend (Hot Reload):**
    ```bash
    cd web && npm run dev
    ```

---

## 🛡 Configuration

Configuration is handled via Environment Variables.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DATABASE_URL` | Postgres connection string | `postgres://user:pass@localhost:5432/opsdeck` |
| `APP_ENV` | `production` or `development` | `development` |
| `AUTH_GOOGLE_CLIENT_ID` | For SSO | - |
| `SMTP_HOST` | For email notifications | - |
| `REDIS_URL` | (Optional) For distributed sessions | - |

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📄 License

This project is licensed under the **AGPLv3 License** - see the [LICENSE](LICENSE) file for details.
