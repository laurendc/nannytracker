import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { vi } from 'vitest'
import Summaries from '../Summaries'

// Mock the API calls
vi.mock('../../lib/api', () => ({
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
const mockSummaries = [
  {
    weekStart: '2024-12-15',
    weekEnd: '2024-12-21',
    totalMiles: 7.3,
    totalAmount: 5.11,
    trips: [
      {
        date: '2024-12-18',
        origin: 'Home',
        destination: 'Work',
        miles: 5.2,
        type: 'single' as const,
      },
      {
        date: '2024-12-19',
        origin: 'Work',
        destination: 'Store',
        miles: 2.1,
        type: 'single' as const,
      },
    ],
    expenses: [
      {
        date: '2024-12-18',
        amount: 15.50,
        description: 'Lunch',
      },
      {
        date: '2024-12-19',
        amount: 8.75,
        description: 'Coffee',
      },
    ],
  },
]

describe('Summaries', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders summaries page title and description', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Weekly Summaries')).toBeInTheDocument()
      expect(screen.getByText(/Overview of your weekly mileage and expense tracking/)).toBeInTheDocument()
    })
  })

  it('displays loading state initially', () => {
    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    // Check for loading skeleton elements
    const skeletonElements = document.querySelectorAll('.animate-pulse')
    expect(skeletonElements.length).toBeGreaterThan(0)
  })

  it('displays weekly summaries when data is loaded', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Weekly Summaries')).toBeInTheDocument()
    })
  })

  it('shows empty state when no summaries exist', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('No weekly summaries available yet.')).toBeInTheDocument()
      expect(screen.getByText(/Add some trips and expenses to see your weekly summaries/)).toBeInTheDocument()
    })
  })

  it('displays summary information correctly when data is available', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText(/Dec \d+ - Dec \d+, 2024/)).toBeInTheDocument()
      // Use getAllByText since there are multiple elements with 7.3
      const elements = screen.getAllByText('7.3')
      expect(elements.length).toBeGreaterThan(0)
      // Use getAllByText since there are multiple elements with $5.11
      const amountElements = screen.getAllByText('$5.11')
      expect(amountElements.length).toBeGreaterThan(0)
    })
  })

  it('shows summary details when expanded', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Trips:')).toBeInTheDocument()
      expect(screen.getByText('Expenses:')).toBeInTheDocument()
    })
  })

  it('displays trip details in summary', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText(/Dec \d+: Home → Work \(5.2 miles\)/)).toBeInTheDocument()
      expect(screen.getByText(/Dec \d+: Work → Store \(2.1 miles\)/)).toBeInTheDocument()
    })
  })

  it('displays expense details in summary', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText(/Dec \d+: Lunch \(\$15.50\)/)).toBeInTheDocument()
      expect(screen.getByText(/Dec \d+: Coffee \(\$8.75\)/)).toBeInTheDocument()
    })
  })

  it('shows correct mileage calculations', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      // Use getAllByText since there are multiple elements with 7.3
      const elements = screen.getAllByText('7.3')
      expect(elements.length).toBeGreaterThan(0)
    })
  })

  it('shows correct expense totals', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('$24.25')).toBeInTheDocument() // 15.50 + 8.75
    })
  })

  it('displays reimbursement calculation', async () => {
    const { summariesApi } = await import('../../lib/api')
    vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

    render(
      <TestWrapper>
        <Summaries />
      </TestWrapper>
    )

    await waitFor(() => {
      // Use getAllByText since there are multiple elements with $5.11
      const elements = screen.getAllByText('$5.11')
      expect(elements.length).toBeGreaterThan(0)
    })
  })
}) 