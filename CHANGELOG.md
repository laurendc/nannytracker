# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.4] - 2024-07-04

### Added
- Optimized build workflow with fast development builds
- New `./scripts/dev-build.sh` for quick development builds with verification
- Enhanced build targets: `build-dev`, `build-ci`, `build-all`
- Comprehensive documentation updates for build workflow
- Build workflow patterns in AI assistant context

### Changed
- CI/CD now uses faster Linux-only builds for PR checks
- Full cross-platform builds reserved for releases only
- Updated contributing guidelines with new build options
- Improved developer onboarding experience
- Enhanced README with build workflow section

### Technical
- Faster development cycles with optimized build targets
- Maintained cross-platform support for release distribution
- Professional CI/CD pipeline optimization
- Industry-standard TUI development practices

## [Unreleased]

### Added
- Version management system
- Release automation with GitHub Actions
- Makefile for common development tasks
- Version endpoints for both TUI and web applications

### Changed
- Improved environment file loading to work from any subdirectory
- Enhanced test coverage for configuration package

### Fixed
- Environment file path resolution when running from subdirectories

## [0.1.0] - 2024-01-XX

### Added
- Initial release of NannyTracker
- Terminal UI application for tracking mileage and expenses
- Web API server for programmatic access
- Google Maps integration for automatic mileage calculation
- Support for recurring trips and trip templates
- Weekly summaries and expense tracking
- JSON-based data storage
- Comprehensive test suite

### Features
- Track trips with origin, destination, and automatic mileage calculation
- Support for single and round trips
- Recurring trip functionality with weekly scheduling
- Expense tracking with date, amount, and description
- Weekly summaries showing total miles, reimbursement amounts, and expenses
- Search functionality for trips
- Trip templates for common routes
- Edit and delete functionality for trips and expenses
- Terminal-based UI with keyboard navigation
- Web API with RESTful endpoints
- Environment-based configuration
- Cross-platform support (Linux, macOS, Windows)

[Unreleased]: https://github.com/laurendc/nannytracker/compare/v1.0.4...HEAD
[1.0.4]: https://github.com/laurendc/nannytracker/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/laurendc/nannytracker/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/laurendc/nannytracker/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/laurendc/nannytracker/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/laurendc/nannytracker/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/laurendc/nannytracker/releases/tag/v0.1.0 