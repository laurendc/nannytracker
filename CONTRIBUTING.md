# Contributing Guidelines for Nanny Tracker

Welcome! This document outlines the steps and conventions for contributing to the Nanny Tracker project. We're glad you're interested in helping out.

## Overview

Nanny Tracker is a terminal-based application (built in Go using the Bubble Tea TUI framework) for tracking mileage, expenses, and reimbursements. It uses Google Maps API for mileage calculations, JSON-based persistent storage, and is configured via environment variables. (See the [README](README.md) for further details.)

## Getting Started

1. Clone the repository:
   • git clone https://github.com/lauren/nannytracker.git
   • cd nannytracker

2. Install dependencies (and ensure you have Go (version 1.23.0 or later) installed):
   • go mod download

3. (Optional) Create a .env file (for example, with your Google Maps API key) and a config.json (if you wish to customize settings) as described in the README.

4. Build the application (for example, "go build") and run tests ("go test ./...").

## Pre-commit Hooks

We use pre-commit (see [.pre-commit-config.yaml](.pre-commit-config.yaml)) to enforce code hygiene and linting. Before you commit, please run (or install) pre-commit so that your changes are automatically checked (for example, trailing whitespace, end-of-file newlines, and Go-specific checks via golangci-lint). (If you're new to pre-commit, please refer to its [documentation](https://pre-commit.com).)

## Commit & Branch Naming Conventions

• Commit messages should be clear and concise. (For example, "feat: add new feature" or "fix: resolve bug in ...".)  
• Branch names should follow a convention such as "feature/..." (for new features), "fix/..." (for bug fixes), "docs/..." (for documentation updates), etc. (See the [git config](.git/config) for examples.)

## Pull Requests (PRs)

• Before submitting a PR, please ensure that your branch is up to date (for example, rebase on main) and that all tests (and pre-commit checks) pass.  
• Please provide a clear title and description (and, if applicable, reference an issue) so that reviewers can understand your change.  

## CI Workflows

Our CI (see [.github/workflows/ci.yml](.github/workflows/ci.yml) and [.github/workflows/security.yml](.github/workflows/security.yml)) runs (among other things) linting (via golangci-lint) and security scans (via gosec, nancy, trufflehog, and goda) on every PR. Please review the output (and fix any issues) that appear.

## License

Nanny Tracker is licensed under the GNU General Public License (GPL) (see [LICENSE](LICENSE)). By contributing, you agree that your contribution is also licensed under the GPL.

---

Thank you for your interest in Nanny Tracker! If you have any questions or need further assistance, please open an issue (or contact the [CODEOWNERS](.github/CODEOWNERS) (currently @laurendc)). 