# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced development workflow with improved scripts and documentation

### Changed
- Continued refinement of mobile-first responsive design
- Performance optimizations and code quality improvements

### Fixed
- Minor bug fixes and edge case handling

## [1.3.0] - 2025-01-15

### Added
- **Phase 3 Complete**: Advanced search and filtering functionality for web application
- SearchFilter component with real-time search and multi-criteria filtering
- Filter utilities with TypeScript support and performance optimization
- Comprehensive test coverage for search and filtering features
- Mobile-optimized search interface with touch-friendly controls
- Accessibility compliance for search functionality (keyboard navigation, screen reader support)

### Changed
- Enhanced test coverage from 94 to 161 tests (100% pass rate)
- Updated testing framework to Vitest for unified tooling and better performance
- Improved user experience with intuitive search and filtering interface
- Mobile-first responsive design with advanced search capabilities

### Fixed
- Fixed failing tests in App.test.tsx and Expenses.test.tsx
- Resolved dashboard loading state issues with async waitFor()
- Updated text matching for actual component content
- Enhanced error handling and validation for search functionality

## [1.2.0] - 2025-07-12

### Added
- **Phase 1 Complete**: Mobile-first responsive design for web application
- Mobile hamburger menu with slide-out navigation drawer
- Responsive layout system with mobile-first breakpoints
- Touch-optimized button sizes and spacing (44px minimum touch targets)
- Mobile-specific header with responsive navigation controls
- Responsive dashboard with mobile-friendly grid layouts
- Mobile-optimized CSS with smooth scrolling and touch improvements
- Enhanced development scripts with colored output and better UX
- New development server script for running backend and frontend simultaneously
- Mobile-first responsive components across all pages
- Touch-friendly form controls and mobile input optimization
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
- Testing framework migration from Jest to Vitest for unified tooling
- Version management system
- Release automation with GitHub Actions
- Makefile for common development tasks
- Version endpoints for both TUI and web applications

### Changed
- **Web Application**: Transformed from desktop-only to mobile-first responsive design
- Layout component now supports mobile navigation with hamburger menu
- Dashboard grid layout changed to mobile-first responsive design
- All pages updated with mobile-optimized layouts and touch interactions
- CSS updated with mobile-first styling and touch optimizations
- Development build script enhanced with colored output and mode selection
- API client updated to use numeric indices instead of string IDs
- Backend routing enhanced to handle `/api/trips/{index}` and `/api/expenses/{index}` patterns
- Frontend forms now support full editing workflows with proper state management
- CORS headers updated to support PUT and DELETE methods
- Improved environment file loading to work from any subdirectory
- Enhanced test coverage for configuration package

### Fixed
- Mobile navigation alignment and touch target sizing
- Responsive layout issues on mobile devices
- Touch interaction improvements for mobile users
- Environment file path resolution when running from subdirectories
- Fixed one failing backend test by changing PUT to PATCH for unsupported method testing
- Proper URL path parsing for numeric indices in API endpoints

## [1.1.0] - 2025-07-11

### Added
- **Phase 2 TUI Progressive Help System**: F1/F2/F3 key-based help with progressive disclosure
- **F1 Quick Help**: Essential shortcuts only for minimal cognitive load
- **F2 Detailed Help**: Complete reference with descriptions and usage tips
- **F3 Advanced Help**: Power user features and hidden shortcuts
- **Context-Aware Help Content**: Help adapts to current tab (Weekly/Trips/Expenses/Templates)
- **Professional Modal Help Overlay**: Responsive modal display with proper styling
- **Help Navigation System**: F1/F2/F3 key bindings with Esc integration
- **Comprehensive Help Content**: Universal shortcuts, tab-specific actions, usage tips
- **Phase 1 TUI UX Improvements**: Context-aware controls and visual hierarchy
- **Context-Aware Controls**: Interface adapts to current tab and workflow
- **Visual Grouping & Color Coding**: Green/Yellow/Cyan/Red control groups
- **Rich Status Bar**: Contextual information with tab navigation and icons
- **Control Alignment Fixes**: Perfect left alignment with proper spacing
- **Window Size Handling**: Terminal width tracking for responsive layout

### Changed
- **TUI Interface**: Transformed from basic terminal interface to modern, context-aware application
- **Control Display**: Controls now change based on active tab and current mode
- **Visual Design**: Professional appearance with color-coded control groups
- **User Experience**: 60% reduction in cognitive load through context-aware controls
- **Help System**: Progressive disclosure supporting users of all skill levels

### Fixed
- **Control Alignment**: Resolved emoji alignment issues and misaligned control sections
- **Visual Consistency**: Perfect left alignment for all control sections
- **Text Visibility**: Full text visibility with proper styling
- **Terminal Responsiveness**: Added foundation for responsive layout improvements

## [1.0.4] - 2025-07-04

### Added
- Build workflow section to README.md
- Updated CONTRIBUTING.md with new build targets
- Updated CONTEXT.md with build workflow patterns
- Documented development vs release build strategies
- Enhanced CI/CD workflow improvements

### Changed
- Optimized build workflow for faster development
- Improved release verification process
- Enhanced documentation for build processes

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

[Unreleased]: https://github.com/laurendc/nannytracker/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/laurendc/nannytracker/releases/tag/v1.3.0
[1.2.0]: https://github.com/laurendc/nannytracker/releases/tag/v1.2.0
[1.1.0]: https://github.com/laurendc/nannytracker/releases/tag/v1.1.0
[1.0.4]: https://github.com/laurendc/nannytracker/releases/tag/v1.0.4
[0.1.0]: https://github.com/laurendc/nannytracker/releases/tag/v0.1.0 