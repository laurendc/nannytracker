# Nanny Mileage Tracker

A simple Terminal User Interface (TUI) application to track nanny mileage and calculate reimbursement. The application uses Google Maps API for accurate distance calculations. NOTE: THIS IS VERY MUCH A WORK IN PROGRESS.

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
2. Enter the origin address and press Enter
3. Enter the destination address and press Enter
4. Enter the trip date in YYYY-MM-DD format and press Enter
5. The application will calculate the actual distance using Google Maps
6. The trip will be saved automatically and added to the history
7. Press Ctrl+C to quit the application

## Data Storage

All trips are automatically saved to `~/.nannytracker/trips.json`. This means your trips will be preserved between sessions.

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
- Add monthly/weekly summaries
- Add date range filtering for trips 