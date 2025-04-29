# Nanny Mileage Tracker

A simple Terminal User Interface (TUI) application to track nanny mileage and calculate reimbursement. The application uses Google Maps API for accurate distance calculations.

## Features

- Enter origin and destination addresses
- Accurate distance calculation using Google Maps Distance Matrix API
- Track multiple trips with dates
- Calculate total mileage and reimbursement
- Simple and intuitive interface
- Automatic trip saving (trips are saved in ~/.nannytracker/trips.json)
- Configurable reimbursement rate (defaults to $0.70 per mile)
- Input validation for addresses, trip data, and dates
- Date tracking for each trip (YYYY-MM-DD format)
- Weekly summaries with total miles and reimbursement amounts
- Persistent storage with automatic data structure creation

## Data Structure

The application stores data in a JSON format with the following structure:

```json
{
  "trips": [
    {
      "origin": "string",
      "destination": "string",
      "miles": number,
      "date": "YYYY-MM-DD"
    }
  ],
  "weekly_summaries": [
    {
      "WeekStart": "YYYY-MM-DD",
      "WeekEnd": "YYYY-MM-DD",
      "TotalMiles": number,
      "TotalAmount": number
    }
  ]
}
```

## Installation

1. Make sure you have Go installed (version 1.23 or higher)
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Set up your environment variables:
   - Create a `.env` file in the project root
   - Add your Google Maps API key: `GOOGLE_MAPS_API_KEY=your_api_key_here`
   - (Optional) Configure custom rate per mile: `RATE_PER_MILE=0.70`
   - (Optional) Configure custom data file path: `DATA_FILE_PATH=~/.nannytracker/trips.json`
5. Run `go run main.go` to start the application

## Usage

1. Launch the application with `go run main.go`
2. Enter the date in YYYY-MM-DD format and press Enter
3. Enter the origin address and press Enter
4. Enter the destination address and press Enter
5. The application will:
   - Calculate the actual distance using Google Maps
   - Save the trip automatically
   - Update weekly summaries
   - Display total mileage and reimbursement
6. Press Ctrl+C to quit the application

## Data Storage

All trips are automatically saved to `~/.nannytracker/trips.json`. The application:
- Creates the data directory if it doesn't exist
- Initializes an empty data structure for new installations
- Maintains weekly summaries automatically
- Preserves all data between sessions

## Configuration

The application can be configured using environment variables:

- `GOOGLE_MAPS_API_KEY` (required): Your Google Maps API key for distance calculations
- `RATE_PER_MILE` (optional): Custom reimbursement rate per mile (default: 0.70)
- `DATA_FILE_PATH` (optional): Custom location for the trips data file

You can set these either in your environment or in a `.env` file in the project root.

## Future Improvements

- Add ability to edit/delete trips
- Add export functionality for reimbursement reports
- Add support for recurring trips
- Add monthly summaries
- Add date range filtering for trips
- Add trip categories or tags
- Add support for multiple reimbursement rates
- Add data backup functionality 