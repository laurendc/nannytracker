# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Phase 2 Complete**: Full CRUD operations for web application
- PUT and DELETE API endpoints for trips and expenses
- Complete edit functionality for trips with all fields (date, origin, destination, miles, type)
- Complete edit functionality for expenses with all fields (date, amount, description)
- Delete operations with user confirmation dialogs
- Index-based item management for consistent API operations
- Comprehensive error handling and loading states
- React Query mutations for optimistic updates
- 10+ new backend tests covering all CRUD operations
- Enhanced frontend test coverage (75 tests passing)
- Version management system
- Release automation with GitHub Actions
- Makefile for common development tasks
- Version endpoints for both TUI and web applications

### Changed
- API client updated to use numeric indices instead of string IDs
- Backend routing enhanced to handle `/api/trips/{index}` and `/api/expenses/{index}` patterns
- Frontend forms now support full editing workflows with proper state management
- CORS headers updated to support PUT and DELETE methods
- Improved environment file loading to work from any subdirectory
- Enhanced test coverage for configuration package

### Fixed
- Environment file path resolution when running from subdirectories
- Fixed one failing backend test by changing PUT to PATCH for unsupported method testing
- Proper URL path parsing for numeric indices in API endpoints

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

[Unreleased]: https://github.com/laurendc/nannytracker/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/laurendc/nannytracker/releases/tag/v0.1.0 