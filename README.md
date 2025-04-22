# Nanny Mileage Tracker

A simple Terminal User Interface (TUI) application to track nanny mileage and calculate reimbursement at $0.70 per mile.

## Features

- Enter origin and destination addresses
- Track multiple trips
- Calculate total mileage and reimbursement
- Simple and intuitive interface
- Automatic trip saving (trips are saved in ~/.nannytracker/trips.json)

## Installation

1. Make sure you have Go installed (version 1.23 or higher)
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Run `go run main.go` to start the application

## Usage

1. Launch the application with `go run main.go`
2. Enter the origin address and press Enter
3. Enter the destination address and press Enter
4. The trip will be saved automatically and added to the history
5. Press Ctrl+C to quit the application

## Data Storage

All trips are automatically saved to `~/.nannytracker/trips.json`. This means your trips will be preserved between sessions.

## Current Limitations

- The mileage is currently hardcoded to 10 miles per trip (this would be replaced with actual Google Maps API integration in a production version)

## Future Improvements

- Add Google Maps API integration for accurate distance calculation
- Add ability to edit/delete trips
- Add date tracking for trips
- Add export functionality for reimbursement reports 