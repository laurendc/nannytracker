# Nanny Tracker

A terminal-based application for tracking mileage, expenses, and calculating reimbursements for nannies.

## Features

- Track trips with date, origin, destination, and mileage
- Track reimbursable expenses with date, amount, and description
- Support for single and round trips
- Support for recurring trips with weekly scheduling
- Automatic mileage calculation using Google Maps API
- **Weekly summaries**: View a summary of your trips and expenses for each week. Use the left (←) and right (→) arrow keys to navigate between weeks.
- **Search trips** by origin, destination, date, or type (Ctrl+F)
- Edit and delete trips with confirmation
- Persistent storage of trip and expense data
- **Beautiful terminal UI** with color highlighting for navigation and clear alignment
- **Ctrl+R**: Toggle recurring trip mode or convert selected trip to recurring

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
- **Ctrl+E**: Edit selected trip or expense
- **Ctrl+D**: Delete selected trip or expense (requires confirmation)
- **Ctrl+X**: Add new expense
- **Ctrl+F**: Toggle search mode (filter trips)
- **↑/↓**: Navigate through trips and expenses (selected item shown in color)
- **Tab**: Switch between trips and expenses list
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
- Add monthly summaries
- Add date range filtering for trips
- Add data backup functionality