# Release Management Guide

This document outlines the release management process for NannyTracker, including versioning strategy, release automation, and best practices.

## Overview

NannyTracker uses a comprehensive release management system that includes:

- **Semantic Versioning (SemVer)** for version numbers
- **Automated release workflow** with GitHub Actions
- **Cross-platform binary distribution**
- **Version tracking** throughout the application
- **Changelog management** for tracking changes

## Versioning Strategy

### Semantic Versioning (SemVer)

We follow the [Semantic Versioning 2.0.0](https://semver.org/) specification:

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes (API changes, major rewrites)
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes and security patches

### Version Examples

- `v1.0.0` - Initial stable release
- `v1.1.0` - New features added
- `v1.1.1` - Bug fixes
- `v2.0.0` - Breaking changes

## Release Process

### 1. Pre-Release Checklist

Before creating a release, ensure:

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Security scan passes (`make security`)
- [ ] CHANGELOG.md is updated
- [ ] Version information is current

### 2. Creating a Release

#### Using Makefile (Recommended)

```bash
# Create a new release
make release VERSION=v1.0.0
```

This command will:
1. Run all tests
2. Build binaries for all platforms
3. Create a git tag
4. Push the tag to trigger GitHub Actions

#### Manual Process

```bash
# 1. Update version in code (if needed)
# 2. Run tests
make test

# 3. Build for all platforms
make build-all

# 4. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 3. Automated Release Workflow

When a tag is pushed, GitHub Actions automatically:

1. **Builds binaries** for multiple platforms:
   - Linux (AMD64, ARM64)
   - macOS (AMD64, ARM64)
   - Windows (AMD64)

2. **Creates a GitHub release** with:
   - Release notes from commits
   - Binary downloads
   - Version information

3. **Uploads assets** to the release

## Version Information

### Application Version

Both the TUI and web applications include version information:

```bash
# TUI application
./nannytracker --version

# Web application

```

### API Version Endpoint

The web server provides version information via API:

```bash
curl http://localhost:8080/version
```

Response:
```json
{
  "version": "1.0.0",
  "build_time": "2024-01-15T10:30:00Z",
  "git_commit": "abc1234",
  "go_version": "go1.23.0",
  "os": "linux",
  "arch": "amd64"
}
```

## Development Workflow

### Daily Development

```bash
# Build for current platform
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Before Committing

```bash
# Run full test suite
make test-race

# Check for security issues
make security

# Update dependencies (if needed)
make deps-update
```

## Dependency Management

### Checking Dependencies

```bash
# Check for outdated dependencies
make deps-check

# Update all dependencies
make deps-update
```

### Dependency Update Policy

1. **Patch updates**: Automatically update
2. **Minor updates**: Review and update monthly
3. **Major updates**: Review carefully, test thoroughly

## Release Notes

### Changelog Structure

The `CHANGELOG.md` file follows the [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
## [Unreleased]

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security fixes
```

### Writing Release Notes

1. **Be descriptive** - Explain what changed and why
2. **Use present tense** - "Adds feature" not "Added feature"
3. **Include breaking changes** - Clearly mark any breaking changes
4. **Reference issues** - Link to relevant GitHub issues/PRs

## Binary Distribution

### Supported Platforms

- **Linux**: AMD64, ARM64
- **macOS**: AMD64, ARM64
- **Windows**: AMD64

### Binary Naming Convention

```
nannytracker-{os}-{arch}[.exe]
```

Examples:
- `nannytracker-linux-amd64`
- `nannytracker-darwin-arm64`
- `nannytracker-windows-amd64.exe`

### Installation Instructions

For each release, include installation instructions:

```bash
# Download appropriate binary
wget https://github.com/laurendc/nannytracker/releases/download/v1.0.0/nannytracker-linux-amd64

# Make executable
chmod +x nannytracker-linux-amd64

# Run
./nannytracker-linux-amd64
```

## Release Schedule

### Release Types

1. **Patch releases** (v1.0.x): Bug fixes and security patches
   - As needed, typically within 1-2 weeks of issues
   
2. **Minor releases** (v1.x.0): New features
   - Monthly or as features are ready
   
3. **Major releases** (vx.0.0): Breaking changes
   - Quarterly or as needed for major changes

### Release Branches

- `main`: Development branch
- `release/v*`: Release branches for major versions
- Tags: Specific release points

## Troubleshooting

### Common Issues

1. **Tag already exists**
   ```bash
   # Remove local tag
   git tag -d v1.0.0
   
   # Remove remote tag
   git push origin --delete v1.0.0
   ```

2. **Build failures**
   ```bash
   # Clean and rebuild
   make clean
   make build-all
   ```

3. **Test failures**
   ```bash
   # Run specific test
   go test ./pkg/config -v
   
   # Run with race detection
   make test-race
   ```

## Best Practices

### Version Management

1. **Always use semantic versioning**
2. **Update CHANGELOG.md before releasing**
3. **Test thoroughly before tagging**
4. **Use descriptive commit messages**

### Release Quality

1. **Run full test suite** before release
2. **Test on multiple platforms** if possible
3. **Verify binary functionality** after build
4. **Update documentation** for new features

### Communication

1. **Announce releases** in appropriate channels
2. **Highlight breaking changes** prominently
3. **Provide migration guides** for major releases
4. **Respond to issues** promptly

## Future Enhancements

Planned improvements to the release management system:

1. **Automated changelog generation** from commit messages
2. **Release candidate workflow** for major releases
3. **Docker image distribution** for containerized deployments
4. **Package manager support** (Homebrew, apt, etc.)
5. **Automated dependency vulnerability scanning** 