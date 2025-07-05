# Contributing to NannyTracker

Thank you for your interest in contributing! This project values consistency, automation, and a great developer experience. Please read these guidelines before submitting code, issues, or documentation.

---

## Development Workflow: Makefile as Source of Truth

NannyTracker uses a **Makefile as the single source of truth** for all build, test, lint, and release logic. This ensures:
- **Consistency** between local development and CI/CD
- **No duplication** of build/test logic in scripts or workflows
- **Easy onboarding** for new contributors

**How it works:**
- All build, test, lint, and release commands are defined in the `Makefile`.
- The CI workflows (GitHub Actions) call the Makefile for all project-specific logic.
- Contributors should use `make` commands locally (e.g., `make test`, `make build`, `make lint`).
- Do **not** duplicate build/test/lint logic in shell scripts or workflow YAML files.

**Benefits:**
- "Works on my machine" = "Works in CI"
- Easy to update build/test logic in one place
- Less risk of drift between local and CI environments

---

## Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/laurendc/nannytracker.git
   cd nannytracker
   ```
2. **Install Go (1.23+) and make sure `make` is available.**
3. **Install dependencies:**
   ```bash
   make deps
   ```
4. **Run tests:**
   ```bash
   make test
   ```
5. **Build the application:**
   ```bash
   # For development (fast, current platform)
   make build
   
   # For quick development builds
   ./scripts/dev-build.sh
   
   # For all platforms (releases only)
   make build-all
   ```
6. **Run linter and security checks:**
   ```bash
   make lint
   make security
   ```

---

## Pull Request Process

1. **Fork the repository and create a feature branch.**
2. **Make your changes.**
3. **Add or update tests for new/changed functionality.**
4. **Run all checks locally:**
   - `make test`
   - `make build` (or `./scripts/dev-build.sh`)
   - `make lint`
   - `make security`
   - `make fmt`
5. **Update `CHANGELOG.md` if your change is user-facing.**
6. **Push your branch and open a Pull Request.**
7. **CI will automatically run all checks using the Makefile.**
8. **Address any review feedback and CI failures.**

---

## Continuous Integration (CI)

- All PRs and pushes to `main` run the full test, lint, and security suite via GitHub Actions.
- The CI workflow calls the Makefile for all project logic.
- Releases are created by tagging a commit and running `make release VERSION=vX.Y.Z`.
- Release artifacts are built and verified automatically.

---

## Questions or Help?

- For build or test issues, always try the Makefile targets first.
- If you find a bug or have a feature request, open an issue.
- For questions about the workflow, ask in your PR or open a discussion.

---

Thank you for helping make NannyTracker better! 