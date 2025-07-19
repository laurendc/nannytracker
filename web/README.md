# NannyTracker Frontend

A modern React + TypeScript frontend for the NannyTracker mileage and expense tracking application.

## Features

- **Dashboard**: Overview of trips, expenses, and weekly summaries
- **Trip Management**: Add, edit, and delete mileage entries
- **Expense Tracking**: Record and manage reimbursable expenses
- **Weekly Summaries**: View detailed weekly reports with charts
- **Advanced Search & Filtering**: Real-time search with multi-criteria filtering
- **Responsive Design**: Mobile-first approach with Tailwind CSS
- **Real-time Data**: React Query for efficient data fetching and caching
- **Complete CRUD Operations**: Full Create, Read, Update, Delete functionality
- **Mobile-First UX**: Touch-optimized controls and responsive navigation

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **React Router** for navigation
- **React Query** for server state management
- **Tailwind CSS** for styling
- **Recharts** for data visualization
- **Vitest** + **React Testing Library** for testing (161 tests passing)
- **MSW** for API mocking

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go backend running on `localhost:8080`

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

### Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## Testing

### Modern Testing with Vitest

This project uses **Vitest** for testing, providing a fast, modern testing experience that integrates seamlessly with Vite.

**Current Test Status:**
- **161 tests passing** (100% pass rate)
- **Backend**: 21 tests passing
- **Frontend**: 140 tests passing
- **Comprehensive coverage** for all features including search and filtering

**Benefits of Vitest:**
- 3-5x faster test execution than Jest
- Unified tooling with Vite (same bundler for dev and testing)
- Hot reload testing with instant feedback
- Interactive test UI for better debugging
- Smaller bundle size and fewer dependencies

### Running Tests

```bash
# Run tests in watch mode
npm test

# Run tests once
npm run test:run

# Run tests with coverage
npm run test:coverage

# Run tests with UI
npm run test:ui
```

### Test Structure

- **Unit Tests**: Test individual components and functions
- **Integration Tests**: Test component interactions and API calls
- **API Tests**: Test API client methods and error handling
- **CRUD Tests**: Test complete Create, Read, Update, Delete workflows

### Test Files

- `src/components/__tests__/` - Component tests (including SearchFilter)
- `src/pages/__tests__/` - Page component tests  
- `src/lib/__tests__/` - API and utility tests
- `src/utils/__tests__/` - Filter utilities and helper tests
- `src/test/` - Test setup and utilities

### Testing Best Practices

- Use React Testing Library for component testing
- Test user interactions, not implementation details
- Mock API calls with MSW (Mock Service Worker)
- Use semantic queries (getByRole, getByLabelText)
- Test accessibility features
- Write tests that resemble how users interact with the app
- Use Vitest's `vi.mock()` for mocking dependencies

## Project Structure

```
web/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── Layout.tsx      # Main layout with navigation
│   │   ├── SearchFilter.tsx # Advanced search and filtering component
│   │   └── __tests__/      # Component tests
│   ├── pages/              # Page components
│   │   ├── Dashboard.tsx   # Dashboard page
│   │   ├── Trips.tsx       # Trip management page
│   │   ├── Expenses.tsx    # Expense tracking page
│   │   ├── Summaries.tsx   # Weekly summaries page
│   │   └── __tests__/      # Page tests
│   ├── lib/                # Utilities and API client
│   │   ├── api.ts          # API client functions
│   │   └── __tests__/      # API tests
│   ├── utils/              # Utility functions
│   │   ├── filterUtils.ts  # Search and filtering logic
│   │   └── __tests__/      # Utility tests
│   ├── types/              # TypeScript type definitions
│   │   └── index.ts        # Shared types
│   ├── test/               # Test setup and utilities
│   │   ├── setup.ts        # Vitest configuration
│   │   ├── utils/          # Test utilities
│   │   └── mocks/          # API mocks
│   ├── App.tsx             # Main app component
│   ├── main.tsx            # App entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
├── dist/                   # Build output
├── package.json            # Dependencies and scripts
├── vite.config.ts          # Vite & Vitest configuration
├── tsconfig.json           # TypeScript configuration
├── tailwind.config.js      # Tailwind CSS configuration
└── README.md               # This file
```

## API Integration

The frontend communicates with the Go backend API running on `localhost:8080`. The API endpoints include:

### Trip Management
- `GET /api/trips` - Get all trips
- `POST /api/trips` - Create a new trip
- `PUT /api/trips/{index}` - Update trip at index
- `DELETE /api/trips/{index}` - Delete trip at index

### Expense Management
- `GET /api/expenses` - Get all expenses
- `POST /api/expenses` - Create a new expense
- `PUT /api/expenses/{index}` - Update expense at index
- `DELETE /api/expenses/{index}` - Delete expense at index

### Reports & Health
- `GET /api/summaries` - Get weekly summaries
- `GET /api/health` - Health check

### API Mocking

In development, API calls are mocked using MSW (Mock Service Worker) to provide a consistent development experience without requiring the backend to be running.

## Development

### Code Style

- Use TypeScript for type safety
- Follow React best practices and hooks
- Use Tailwind CSS for styling
- Write meaningful component and function names
- Add JSDoc comments for complex functions

### State Management

- Use React Query for server state
- Use React state for local component state
- Avoid prop drilling with proper component composition

### Performance

- Use React.memo for expensive components
- Implement proper loading states
- Use React Query's caching capabilities
- Optimize bundle size with code splitting

### Testing Workflow

- Write tests for new features using Vitest
- Use `npm test` for watch mode during development
- Run `npm run test:coverage` to check coverage
- Use `npm run test:ui` for interactive debugging

## Deployment

The frontend can be deployed to any static hosting service:

1. Build the project: `npm run build`
2. Deploy the `dist/` directory to your hosting service
3. Configure the API endpoint for production

### Environment Variables

- `VITE_API_URL` - Backend API URL (defaults to `http://localhost:8080`)

## Contributing

1. Write tests for new features using Vitest
2. Ensure all tests pass with `npm run test:run`
3. Follow the existing code style
4. Update documentation as needed

## Troubleshooting

### Common Issues

- **API Connection Errors**: Ensure the Go backend is running on port 8080
- **Build Errors**: Check TypeScript types and dependencies
- **Test Failures**: Verify API mocks are working correctly with Vitest

### Development Tips

- Use the React DevTools for debugging
- Check the browser console for errors
- Use the Network tab to debug API calls
- Use the React Query DevTools for state debugging
- Use Vitest UI (`npm run test:ui`) for interactive test debugging
- Run tests in watch mode for instant feedback during development 