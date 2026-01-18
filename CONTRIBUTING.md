# Contributing to OpsDeck

First off, thank you for considering contributing to OpsDeck! It's people like you that make the open-source community such an amazing place to learn, inspire, and create.

Following these guidelines helps communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

## 📜 Code of Conduct

This project and everyone participating in it is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## 🛠 Getting Started

### Prerequisites
* **Go**: v1.24 or higher
* **Node.js**: v20 or higher
* **Docker**: Required for running the local database and Redis.
* **Make**: (Optional) For running shorthand commands.

### Local Development Environment
1.  **Fork** the repository on GitHub.
2.  **Clone** your fork locally:
    ```bash
    git clone [https://github.com/your-username/opsdeck.git](https://github.com/your-username/opsdeck.git)
    cd opsdeck
    ```
3.  **Start Infrastructure** (Postgres & Redis):
    ```bash
    docker compose -f docker-compose.dev.yml up -d
    ```
4.  **Run Backend**:
    ```bash
    go run cmd/server/main.go
    ```
5.  **Run Frontend** (in a separate terminal):
    ```bash
    cd web && npm install && npm run dev
    ```

---

## 🐛 Found a Bug?

If you find a bug in the source code, you can help us by [submitting an issue](https://github.com/username/opsdeck/issues) to our GitHub Repository. Even better, you can submit a Pull Request with a fix.

**Please include:**
1.  Your OS and Browser version.
2.  Steps to reproduce the bug.
3.  Expected vs. Actual behavior.
4.  Screenshots or logs if applicable.

---

## 💡 Submitting a Pull Request

The core team monitors Pull Requests closely. We will review your PR and either merge it, request changes, or close it with an explanation.

### Process
1.  **Create a Branch:** Create a new branch for your feature or fix.
    * `feat/add-asset-tracking`
    * `fix/login-timeout`
2.  **Coding Standards:**
    * **Go:** We strictly follow `gofmt`. Please run `go fmt ./...` before committing.
    * **React:** We use Prettier. Please run `npm run format` in the `web/` directory.
3.  **Tests:**
    * New features must include unit tests.
    * Bug fixes must include a regression test that fails without the fix and passes with it.
4.  **Commit Messages:** Please follow [Conventional Commits](https://www.conventionalcommits.org/):
    * `feat: add QR code scanning for assets`
    * `fix(ui): correct z-index on mobile modal`

### Developer Certificate of Origin (DCO)

To protect the project's intellectual property and ensure we have the legal right to distribute your code, all PRs must include a "Signed-off-by" line in the commit message.

```text
Signed-off-by: Jane Doe <jane.doe@example.com>

```

You can do this automatically by adding the `-s` flag to your git commit:

```bash
git commit -s -m "feat: my new feature"

```

By signing off, you certify that you wrote the code or have the right to contribute it under the project's license (AGPLv3).

---

## 🎨 Style Guides

### Backend (Go)

* **Error Handling:** Don't ignore errors. Wrap them with context using `fmt.Errorf("doing x: %w", err)`.
* **Concurrency:** Use channels for communication, mutexes for state. Avoid global state where possible.
* **Structure:** Keep business logic in `internal/core`, infrastructure in `internal/adapters`.

### Frontend (React + TS)

* **Components:** Functional components with Hooks only.
* **State:** Use React Context for global state only when necessary. Prefer local state.
* **Styling:** We use Tailwind CSS. Avoid writing custom CSS files unless absolutely necessary.

## ❓ Questions?

Feel free to open a "Discussion" on GitHub if you have questions about the architecture or roadmap before starting work.
