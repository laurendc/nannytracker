# Nanny Tracker

A comprehensive application for tracking mileage, expenses, and calculating reimbursements. Originally built to provide nannies with weekly reimbursements, it can be used by anyone who needs to track expenses and mileage for work purposes.

## Current Status

**Terminal Application**: ✅ **Production Ready** - A fully functional terminal-based application with a rich TUI interface.

**Web Application**: ✅ **Phase 3 Complete** - A modern React-based web interface with **Full CRUD Operations** and **Advanced Search & Filtering**. Mobile-first responsive design with comprehensive functionality.

## Features

### Terminal Application (Production Ready)
- **Rich TUI Interface**: Terminal-based user interface with keyboard navigation
- **Trip Management**: Track trips with date, origin, destination, and automatic mileage calculation
- **Expense Tracking**: Record reimbursable expenses with date, amount, and description
- **Trip Templates**: Create reusable templates for common trips
- **Recurring Trips**: Set up weekly recurring trips with automatic generation
- **Weekly Summaries**: View detailed weekly reports with itemized trips and expenses
- **Search & Filter**: Real-time search through trips and expenses
- **Data Validation**: Comprehensive validation for all entries
- **Persistent Storage**: JSON-based data storage with backup capabilities

### Web Application (Phase 3 Complete - Full CRUD Operations & Advanced Features)
- **Modern React Interface**: Built with React 18, TypeScript, Tailwind CSS
- **Complete CRUD Operations**: ✅ Create, Read, Update, Delete for trips and expenses
- **Advanced Search & Filtering**: Real-time search with multi-criteria filtering
- **Real-time Data**: React Query for efficient data fetching and caching
- **Dashboard**: Overview of trips, expenses, and weekly summaries
- **Interactive Forms**: Full editing capabilities with validation and error handling
- **Confirmation Dialogs**: User-friendly delete confirmations with proper UX
- **API Integration**: RESTful API backend with full CRUD endpoint support
- **Mobile-First Design**: Responsive design with touch-optimized controls

**Phase 3 Features Implemented:**
- ✅ Advanced search and filtering with debounced real-time search
- ✅ Multi-criteria filtering (date range, amount, description, type)
- ✅ Mobile-first responsive design with touch-optimized controls
- ✅ Comprehensive test coverage with 161 tests passing (100% pass rate)
- ✅ Enhanced user experience with intuitive search interface
- ✅ Accessibility compliance with keyboard navigation and screen reader support

**Phase 2 Features (Previously Implemented):**
- ✅ Trip editing with all fields (date, origin, destination, miles, type)
- ✅ Expense editing with all fields (date, amount, description)
- ✅ Delete operations with confirmation dialogs
- ✅ PUT/DELETE API endpoints for trips and expenses
- ✅ Index-based item management for consistent operations
- ✅ Optimistic updates with React Query mutations
- ✅ Comprehensive error handling and user feedback

**Next Phase - Performance & Advanced Features:**
- 🚧 Bundle optimization and code splitting
- 🚧 Progressive Web App features
- 🚧 Advanced data visualization
- 🚧 Export capabilities (PDF/CSV)

## Technical Architecture

### Backend (Go)
- **Core Logic**: Shared business logic between TUI and web applications
- **REST API**: HTTP server providing JSON endpoints for web frontend
- **File Storage**: JSON-based persistent storage with data validation
- **Configuration**: Environment-based configuration with `.env` support
- **Cross-platform**: Supports Linux and macOS

### Frontend (React/TypeScript)
- **Modern Stack**: React 18, TypeScript, Vite, Tailwind CSS
- **State Management**: React Query for server state, React state for local state
- **Testing**: Comprehensive test suite with Vitest and React Testing Library (161 tests passing)
- **Development**: Hot reload, linting, and type checking

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Visit the [releases page](https://github.com/laurendc/nannytracker/releases)
2. Download the appropriate binary for your platform:
   - **Linux**: `nannytracker-linux-amd64` or `nannytracker-linux-arm64`
   - **macOS**: `nannytracker-darwin-amd64` or `nannytracker-darwin-arm64`

3. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x nannytracker-linux-amd64
   ```

4. Run the application:
   ```bash
   # Terminal application
   ./nannytracker-linux-amd64
   ```

### Option 2: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/laurendc/nannytracker.git
   cd nannytracker
   ```

2. Install dependencies:
   ```bash
   # Go dependencies
   go mod download
   
   # Web frontend dependencies (optional)
   cd web && npm install
   ```

3. Build the application:
   ```bash
   # Build terminal application
   go build -o nannytracker ./cmd/tui
   
   # Or use the Makefile
   make build
   ```

## Configuration

1. Create a `.env` file in the project root with your Google Maps API key:
   ```
   GOOGLE_MAPS_API_KEY=your_api_key_here
   ```

2. (Optional) Create a `config.json` file to customize settings:
   ```json
   {
     "rate_per_mile": 0.70,
     "data_path": "~/.nannytracker"
   }
   ```

## Usage

### Terminal Application

```bash
# Run the TUI application
./nannytracker

# Check version information
./nannytracker --version
```

**Keyboard Controls:**
- **Enter**: Confirm input or move to next field
- **Ctrl+E**: Edit selected item
- **Ctrl+D**: Delete selected item (requires confirmation)
- **Ctrl+X**: Add new expense
- **Ctrl+F**: Toggle search mode
- **Ctrl+T**: Create new trip template
- **Ctrl+U**: Use selected template to create a new trip
- **↑/↓**: Navigate through items
- **Tab/Shift+Tab**: Switch between tabs
- **Ctrl+C**: Quit application

### Web Application

The web application is designed to run as a development server or be deployed as a web service, not as a standalone binary.

**For Development:**
```bash
# Start the React development server
cd web && npm run dev

# Build for production
cd web && npm run build
```

### Web API

The web server provides a comprehensive REST API with full CRUD operations:

```bash
# Health check
curl http://localhost:8080/health

# Version information
curl http://localhost:8080/version

# Trips API (Full CRUD)
curl http://localhost:8080/api/trips                                    # GET all trips
curl -X POST http://localhost:8080/api/trips -d '{"date":"2024-01-01","origin":"Home","destination":"Work","miles":10,"type":"single"}' # CREATE
curl -X PUT http://localhost:8080/api/trips/0 -d '{"date":"2024-01-02","origin":"Home","destination":"Work","miles":12,"type":"single"}' # UPDATE
curl -X DELETE http://localhost:8080/api/trips/0                        # DELETE

# Expenses API (Full CRUD)
curl http://localhost:8080/api/expenses                                  # GET all expenses
curl -X POST http://localhost:8080/api/expenses -d '{"date":"2024-01-01","amount":25.50,"description":"Gas"}' # CREATE
curl -X PUT http://localhost:8080/api/expenses/0 -d '{"date":"2024-01-02","amount":30.00,"description":"Parking"}' # UPDATE
curl -X DELETE http://localhost:8080/api/expenses/0                     # DELETE

# Weekly summaries (read-only)
curl http://localhost:8080/api/summaries                                # GET summaries
```

**API Endpoints:**
- `GET /api/trips` - List all trips
- `POST /api/trips` - Create a new trip
- `PUT /api/trips/{index}` - Update trip at index
- `DELETE /api/trips/{index}` - Delete trip at index
- `GET /api/expenses` - List all expenses
- `POST /api/expenses` - Create a new expense
- `PUT /api/expenses/{index}` - Update expense at index
- `DELETE /api/expenses/{index}` - Delete expense at index
- `GET /api/summaries` - Get weekly summaries (read-only)

## Development

### Quick Start

```bash
# Clone the repository
git clone https://github.com/laurendc/nannytracker.git
cd nannytracker

# Install development dependencies
make deps

# Run tests
make test

# Build for development (fast, current platform)
make build

# Quick development build with verification
./scripts/dev-build.sh

# Build for all platforms (releases only)
make build-all

# Run linter
make lint

# Format code
make fmt
```

### Build Workflow

NannyTracker uses an optimized build workflow to balance development speed with cross-platform support:

```bash
# Development builds (fast, current platform)
make build              # Standard development build
./scripts/dev-build.sh  # Quick build with verification

# Release builds (all platforms)
make build-all          # Full cross-platform build for releases
```

**Build Strategy:**
- **Development**: Fast Linux builds for daily work
- **CI/CD**: Linux-only builds for PR checks (faster)
- **Releases**: Full cross-platform builds (Linux, macOS)

### Running Tests

```bash
# Run all Go tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage
make test-coverage

# Run web frontend tests
cd web && npm test
```

### Release Management

NannyTracker uses a comprehensive release management system with automated builds and versioning.

#### Creating a Release

```bash
# Create a new release (requires VERSION=)
make release VERSION=v1.0.0
```

This will:
1. Run all tests
2. Build binaries for all platforms
3. Create a git tag
4. Trigger GitHub Actions to create a release

#### Version Information

```bash
# Check version from command line
./nannytracker --version

# Check version via API
curl http://localhost:8080/version
```

For detailed release management information, see [docs/RELEASE_MANAGEMENT.md](docs/RELEASE_MANAGEMENT.md).

### Project Structure

```
.
├── cmd/
│   ├── tui/
│   │   └── main.go      # Terminal application entry point
│   └── web/
│       └── main.go      # Web API server entry point
├── internal/
│   └── tui/             # Terminal UI components
├── pkg/
│   ├── config/          # Configuration management
│   ├── core/            # Core business logic
│   └── version/         # Version management
├── web/                 # React web frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── lib/         # Utilities and API client
│   │   └── types/       # TypeScript types
│   ├── package.json     # Frontend dependencies
│   └── README.md        # Frontend documentation
├── docs/                # Documentation
├── Makefile             # Development tasks
├── CHANGELOG.md         # Release history
└── README.md            # This file
```

### Dependencies

**Backend (Go):**
- github.com/charmbracelet/bubbletea - Terminal UI framework
- github.com/joho/godotenv - Environment configuration
- Google Maps API - Mileage calculations

**Frontend (React):**
- React 18 with TypeScript
- Vite for development and building
- React Router for navigation
- React Query for server state management
- Tailwind CSS for styling
- Recharts for data visualization

### Development Workflow

1. **Make changes** in a feature branch
2. **Run tests** to ensure everything works
3. **Update CHANGELOG.md** with your changes
4. **Create a release** when ready
5. **Monitor feedback** and iterate

## Version History

See [CHANGELOG.md](CHANGELOG.md) for a complete history of changes and releases.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

For detailed contributing guidelines, see [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md).

## Release Management

For detailed release management information, see [docs/RELEASE_MANAGEMENT.md](docs/RELEASE_MANAGEMENT.md).

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Roadmap

### Current Phase - Mobile-First Implementation (Phase 1)
- **Responsive Design**: Mobile-first approach with touch-optimized controls
- **Navigation**: Mobile-friendly navigation patterns (bottom tabs, slide-out menus)
- **Touch Interface**: Gesture-based interactions and mobile UX patterns
- **Performance**: Optimize for mobile device constraints

### Short Term
- Complete Phase 1 mobile-first implementation
- Add export functionality for reimbursement reports (PDF/CSV)
- Implement trip templates and recurring trips
- Add advanced search and filtering capabilities
- Bundle size optimization (currently 677KB)

### Long Term
- Progressive Web App features (offline support, installation)
- Cloud synchronization
- Multi-user support
- Advanced reporting and analytics
- Data visualization enhancements