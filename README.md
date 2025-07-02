# Nanny Tracker

A comprehensive application for tracking mileage, expenses, and calculating reimbursements. Originally built to provide nannies with weekly reimbursements, it can be used by anyone who needs to track expenses and mileage for work purposes.

## Current Status

**Terminal Application**: ✅ **Production Ready** - A fully functional terminal-based application with a rich TUI interface.

**Web Application**: 🚧 **In Development** - A modern React-based web interface is currently being developed alongside the terminal application.

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

### Web Application (In Development)
- **Modern React Interface**: Built with React 18, TypeScript, and Tailwind CSS
- **Responsive Design**: Mobile-first approach with beautiful UI
- **Real-time Data**: React Query for efficient data fetching and caching
- **Dashboard**: Overview of trips, expenses, and weekly summaries
- **Interactive Charts**: Data visualization with Recharts
- **API Integration**: RESTful API backend for programmatic access

## Technical Architecture

### Backend (Go)
- **Core Logic**: Shared business logic between TUI and web applications
- **REST API**: HTTP server providing JSON endpoints for web frontend
- **File Storage**: JSON-based persistent storage with data validation
- **Configuration**: Environment-based configuration with `.env` support
- **Cross-platform**: Supports Linux, macOS, and Windows

### Frontend (React/TypeScript)
- **Modern Stack**: React 18, TypeScript, Vite, Tailwind CSS
- **State Management**: React Query for server state, React state for local state
- **Testing**: Comprehensive test suite with Jest and React Testing Library
- **Development**: Hot reload, linting, and type checking

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Visit the [releases page](https://github.com/laurendc/nannytracker/releases)
2. Download the appropriate binary for your platform:
   - **Linux**: `nannytracker-linux-amd64` or `nannytracker-linux-arm64`
   - **macOS**: `nannytracker-darwin-amd64` or `nannytracker-darwin-arm64`
   - **Windows**: `nannytracker-windows-amd64.exe`

3. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x nannytracker-linux-amd64
   ```

4. Run the application:
   ```bash
   # Terminal application
   ./nannytracker-linux-amd64
   
   # Web server
   ./nannytracker-web-linux-amd64
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

3. Build the applications:
   ```bash
   # Build terminal application
   go build -o nannytracker ./cmd/tui
   
   # Build web server
   go build -o nannytracker-web ./cmd/web
   
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

```bash
# Start the web server
./nannytracker-web

# Check version information
./nannytracker-web --version
```

**Access the Web Interface:**
1. Start the web server: `./nannytracker-web`
2. Open your browser to `http://localhost:8080`
3. The web interface will be available at the root URL

**For Development:**
```bash
# Start the React development server
cd web && npm run dev

# Build for production
cd web && npm run build
```

### Web API

The web server provides a REST API for programmatic access:

```bash
# Health check
curl http://localhost:8080/health

# Version information
curl http://localhost:8080/version

# Get trips
curl http://localhost:8080/api/trips

# Get expenses
curl http://localhost:8080/api/expenses

# Get weekly summaries
curl http://localhost:8080/api/summaries
```

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

# Build for current platform
make build

# Build for all platforms
make build-all

# Run linter
make lint

# Format code
make fmt
```

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

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Roadmap

### Short Term
- Complete web frontend development
- Add export functionality for reimbursement reports
- Add monthly summaries
- Add date range filtering for trips

### Long Term
- Mobile-friendly web interface
- Cloud synchronization
- Multi-user support
- Data backup functionality
- Advanced reporting and analytics