# Nanny Tracker

A terminal-based application for tracking mileage, expenses, and calculating reimbursements. While this was originally built to provide my nanny with weekly reumbusements, it can also be used for tracking expenses and mileage for anyone.

## Technical Details

- Built in Go using the Bubble Tea TUI framework
- Uses Google Maps API for mileage calculations
- JSON-based persistent storage
- Environment-based configuration
- Comprehensive test coverage

## Features

- Track trips with date, origin, destination, and mileage
- Track reimbursable expenses with date, amount, and description
- Support for single and round trips
- Support for recurring trips with weekly scheduling
- Automatic mileage calculation using Google Maps API
- Weekly summaries of mileage, expenses, and reimbursement amounts
- Search trips by origin, destination, date, or type (Ctrl+F)
- Edit and delete trips with confirmation
- Persistent storage of trip and expense data
- Data validation for all entries
- Automatic trip generation from recurring trips
- Support for custom mileage reimbursement rates
- Configurable data storage location
- Terminal-based UI with keyboard navigation
- Real-time search filtering
- Weekly summaries with itemized trips and expenses
- Manage reusable trip templates for common trips

## Data Structures

The application uses the following core data structures:

- **Trip**: Represents a single trip with origin, destination, mileage, date, and type (single/round)
- **RecurringTrip**: Represents a weekly recurring trip with start/end dates and weekday
- **Expense**: Represents a reimbursable expense with date, amount, and description
- **WeeklySummary**: Contains aggregated data for a week including total miles, reimbursement amount, and expenses
- **TripTemplate**: Represents a reusable template for trips, including name, origin, destination, type, and notes

## Installation

1. Clone the repository:
```bash
git clone https://github.com/laurendc/nannytracker.git
cd nannytracker
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o nannytracker ./cmd/tui
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
- **Ctrl+E**: Edit selected trip, expense, or template
- **Ctrl+D**: Delete selected trip, expense, or template (requires confirmation)
- **Ctrl+X**: Add new expense
- **Ctrl+F**: Toggle search mode (filter trips)
- **Ctrl+T**: Create new trip template
- **Ctrl+U**: Use selected template to create a new trip
- **↑/↓**: Navigate through trips, expenses, or templates (selected item shown in color)
- **Tab/Shift+Tab**: Switch between Weekly Summaries, Trips, Expenses, and Trip Templates tabs
- **Ctrl+C**: Quit application

### Searching Trips

- Press **Ctrl+F** to enter search mode
- Type a search term (origin, destination, date, or type)
- The trip list will be filtered in real time
- Press **Ctrl+F** again to exit search mode and return to the full list

### Adding a Trip

1. Enter the date (YYYY-MM-DD)
2. Enter the origin address
3. Enter the destination address
4. Select trip type (single or round)
5. The mileage will be automatically calculated (doubled for round trips)

### Adding a Recurring Trip

1. Press Ctrl+R to enter recurring trip mode
2. Enter the start date (YYYY-MM-DD)
3. Enter the weekday (0-6, where 0 is Sunday)
4. Enter the origin address
5. Enter the destination address
6. Select trip type (single or round)
7. The mileage will be automatically calculated (doubled for round trips)
8. Trips will be generated for each occurrence of the weekday until the end of the current month

### Adding an Expense

1. Press Ctrl+X to enter expense mode
2. Enter the date (YYYY-MM-DD)
3. Enter the expense amount
4. Enter a description of the expense
5. The expense will be added to the weekly summary for that date

### Using Trip Templates

#### Creating a Template
1. Press **Ctrl+T** from the main screen (date mode)
2. Enter a template name
3. Enter the origin address
4. Enter the destination address
5. Enter the trip type (single or round)
6. (Optional) Enter notes for the template
7. The template will be saved for future use

#### Navigating and Managing Templates
1. Press **Tab** or **Shift+Tab** to switch to the Trip Templates tab
2. Use **↑/↓** to select a template
3. Press **Ctrl+E** to edit the selected template
4. Press **Ctrl+D** to delete the selected template (confirmation required)

#### Using Templates to Create Trips
1. Press **Tab** or **Shift+Tab** to switch to the Trip Templates tab
2. Use **↑/↓** to select the template you want to use
3. Press **Ctrl+U** to create a new trip from the template
4. Enter the date for the new trip
5. The trip will be created with the template's origin, destination, and type
6. The mileage will be automatically calculated based on the origin and destination

### Editing a Trip

1. Select the trip using ↑/↓ keys
2. Press Ctrl+E to enter edit mode
3. Update the fields you want to change:
   - Press Enter without typing to keep the existing value
   - Type a new value and press Enter to update the field
4. The mileage will be automatically recalculated if origin or destination changes

### Editing an Expense

1. Select the expense using ↑/↓ keys (use Tab to switch between trips and expenses)
2. Press Ctrl+E to enter edit mode
3. Update the fields you want to change:
   - Press Enter without typing to keep the existing value
   - Type a new value and press Enter to update the field

### Deleting a Trip

1. Select the trip using ↑/↓ keys
2. Press Ctrl+D to enter delete confirmation mode
3. Type 'yes' and press Enter to confirm deletion
4. Type anything else and press Enter to cancel

### Deleting an Expense

1. Select the expense using ↑/↓ keys (use Tab to switch between trips and expenses)
2. Press Ctrl+D to enter delete confirmation mode
3. Type 'yes' and press Enter to confirm deletion
4. Type anything else and press Enter to cancel

## Weekly Summaries

- Weekly summaries are always displayed at the top of the UI
- Each summary shows total miles, mileage amount, and expenses for the week
- Formatting and alignment have been improved for clarity

## Recurring Trips

- Recurring trips are displayed in a separate section above the regular trips
- Each recurring trip shows the origin, destination, mileage, type, and weekday
- Generated trips from recurring trips appear in the regular trips list
- Recurring trips automatically generate new trips until the end of the current month
- Converting a trip to recurring will remove the original trip and create a new recurring trip

## Development

### Running Tests

All tests have been updated to match the latest UI and logic. To run the full suite:

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/
│   └── tui/
│       └── main.go      # TUI application entry point
├── internal/
│   ├── tui/             # Terminal UI components
│   ├── storage/         # Data persistence
│   └── maps/            # Google Maps integration
├── pkg/                 # Public packages (core logic, etc.)
├── web/                 # Web frontend (future)
└── README.md            # Documentation
```

### Dependencies

- github.com/charmbracelet/bubbletea - Terminal UI framework
- github.com/joho/godotenv - Environment configuration
- Google Maps API - Mileage calculations

## Future Enhancements

- Add export functionality for reimbursement reports
- Add monthly summaries
- Add date range filtering for trips
- Add data backup functionality