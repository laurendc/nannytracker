# Nanny Mileage Tracker

A terminal-based application for tracking mileage and calculating reimbursements for nannies.

## Features

- Track trips with date, origin, destination, and mileage
- Automatic mileage calculation using Google Maps API
- Weekly summaries of mileage and reimbursement amounts
- Edit and delete trips with confirmation
- Persistent storage of trip data
- Beautiful terminal UI with keyboard navigation

## Installation

1. Clone the repository:
```bash
git clone https://github.com/lauren/nannytracker.git
cd nannytracker
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build
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

Run the application:
```bash
./nannytracker
```

### Keyboard Controls

- **Enter**: Confirm input or move to next field
- **Ctrl+E**: Edit selected trip
- **Ctrl+D**: Delete selected trip (requires confirmation)
- **↑/↓**: Navigate through trips (selected trip shown in theme color)
- **Ctrl+C**: Quit application

### Adding a Trip

1. Enter the date (YYYY-MM-DD)
2. Enter the origin address
3. Enter the destination address
4. The mileage will be automatically calculated

### Editing a Trip

1. Select the trip using ↑/↓ keys
2. Press Ctrl+E to enter edit mode
3. Update the date, origin, and destination
4. Press Enter after each field

### Deleting a Trip

1. Select the trip using ↑/↓ keys
2. Press Ctrl+D to enter delete confirmation mode
3. Type 'yes' and press Enter to confirm deletion
4. Type anything else and press Enter to cancel

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

- `cmd/`: Main application entry point
- `internal/`: Core application code
  - `maps/`: Google Maps API integration
  - `model/`: Data models and business logic
  - `storage/`: Data persistence
  - `ui/`: Terminal user interface
- `pkg/`: Shared utilities
  - `config/`: Configuration management

## Future Enhancements

- Add export functionality for reimbursement reports
- Add support for recurring trips
- Add monthly summaries
- Add date range filtering for trips
- Add data backup functionality 