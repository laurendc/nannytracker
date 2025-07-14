import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { vi } from 'vitest'
import Dashboard from '../Dashboard'

// Mock date-fns format function
vi.mock('date-fns', () => ({
  format: vi.fn((date, formatStr) => {
    // Return predictable date strings for testing
    if (formatStr === 'MMM d') {
      return 'Dec 15'
    }
    if (formatStr === 'MMM d, yyyy') {
      return 'Dec 15, 2024'
    }
    return 'Dec 15'
  }),
}))

// Mock the API calls
vi.mock('../../lib/api', () => ({
  tripsApi: {
    getAll: vi.fn(),
  },
  expensesApi: {
    getAll: vi.fn(),
  },
  summariesApi: {
    getAll: vi.fn(),
  },
}))

// Create a simple wrapper for testing
const TestWrapper = ({ children }: { children: React.ReactNode }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {children}
      </BrowserRouter>
    </QueryClientProvider>
  )
}

// Mock data
const mockTrips = [
  {
    date: '2024-12-19',
    origin: 'Home',
    destination: 'Work',
    type: 'single' as const,
    miles: 15.5,
  },
  {
    date: '2024-12-20',
    origin: 'Work',
    destination: 'Store',
    type: 'round' as const,
    miles: 8.2,
  },
]

const mockExpenses = [
  {
    date: '2024-12-19',
    description: 'Lunch',
    amount: 12.5,
  },
  {
    date: '2024-12-20',
    description: 'Coffee',
    amount: 3.75,
  },
]

const mockSummaries = [
  {
    weekStart: '2024-12-15',
    weekEnd: '2024-12-21',
    totalMiles: 23.7,
    totalAmount: 16.25,
    expenses: mockExpenses,
    trips: mockTrips,
  },
]

describe('Dashboard', () => {
  beforeEach(async () => {
    const { tripsApi, expensesApi, summariesApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)
  })

  it('renders dashboard title and description', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
      expect(screen.getByText(/Overview of your mileage and expense tracking/)).toBeInTheDocument()
    })
  })

  it('displays loading state initially', () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    // Check for loading skeleton elements
    const skeletonElements = document.querySelectorAll('.animate-pulse')
    expect(skeletonElements.length).toBeGreaterThan(0)
  })

  it('displays stats cards when data is loaded', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getAllByText('Total Trips').length).toBeGreaterThan(0)
      expect(screen.getAllByText('Total Miles').length).toBeGreaterThan(0)
      expect(screen.getAllByText('Total Expenses').length).toBeGreaterThan(0)
      expect(screen.getAllByText('Reimbursement').length).toBeGreaterThan(0)
      expect(screen.getByText('2')).toBeInTheDocument() // trips.length
      expect(screen.getAllByText('23.7').length).toBeGreaterThan(0) // total miles
      expect(screen.getAllByText('$16.25').length).toBeGreaterThan(0) // total expenses
      expect(screen.getAllByText('$16.59').length).toBeGreaterThan(0) // reimbursement (23.7 * 0.7)
    })
  })

  it('displays recent trips section', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Recent Trips')).toBeInTheDocument()
    })
  })

  it('displays recent expenses section', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Recent Expenses')).toBeInTheDocument()
    })
  })

  it('displays current week summary when available', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Current Week Summary')).toBeInTheDocument()
      // Check for the mocked date format
      expect(screen.getByText('Dec 15 - Dec 15, 2024')).toBeInTheDocument()
      expect(screen.getAllByText('23.7').length).toBeGreaterThan(0)
      expect(screen.getAllByText('$16.25').length).toBeGreaterThan(0)
    })
  })

  it('shows correct trip data', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Home → Work')).toBeInTheDocument()
      expect(screen.getByText('Work → Store')).toBeInTheDocument()
    })
  })

  it('shows correct expense data', async () => {
    render(
      <TestWrapper>
        <Dashboard />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Lunch')).toBeInTheDocument()
      expect(screen.getByText('Coffee')).toBeInTheDocument()
    })
  })
}) 