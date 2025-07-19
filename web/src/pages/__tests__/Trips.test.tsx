import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { vi } from 'vitest'
import Trips from '../Trips'

// Mock the API calls
vi.mock('../../lib/api', () => ({
  tripsApi: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
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
    date: '2024-12-20',
    origin: 'Home',
    destination: 'Work',
    type: 'single' as const,
    miles: 15.5,
  },
  {
    date: '2024-12-21',
    origin: 'Work',
    destination: 'Store',
    type: 'round' as const,
    miles: 8.2,
  },
  {
    date: '2024-12-22',
    origin: 'Home',
    destination: 'Gym',
    type: 'single' as const,
    miles: 5.0,
  },
]

describe('Trips', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Clear localStorage before each test to ensure clean state
    localStorage.clear()
  })

  it('renders trips page title and description', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Trips')).toBeInTheDocument()
      expect(screen.getByText(/Manage your mileage tracking entries/)).toBeInTheDocument()
    })
  })

  it('shows add trip button', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Trip')).toBeInTheDocument()
    })
  })

  it('opens add trip form when button is clicked', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Trip')).toBeInTheDocument()
    })

    const addButton = screen.getByText('Add Trip')
    await act(async () => {
      await user.click(addButton)
    })

    expect(screen.getByText('Add New Trip')).toBeInTheDocument()
    const emptyInputs = screen.getAllByDisplayValue('')
    expect(emptyInputs.length).toBeGreaterThan(0)
    expect(screen.getByRole('combobox')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Enter origin address')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Enter destination address')).toBeInTheDocument()
  })

  it('allows form input and submission', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])
    vi.mocked(tripsApi.create).mockResolvedValue({
      date: '2024-12-20',
      origin: 'Home',
      destination: 'Work',
      type: 'single',
      miles: 0
    })

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Trip')).toBeInTheDocument()
    })

    const addButton = screen.getByText('Add Trip')
    await user.click(addButton)

    // Wait for form to appear
    await waitFor(() => {
      expect(screen.getByText('Add New Trip')).toBeInTheDocument()
    })

    // Select the date input by filtering all empty inputs for type="date"
    const emptyInputs = screen.getAllByDisplayValue('')
    const dateInput = emptyInputs.find(input => input.getAttribute('type') === 'date')
    const originInput = screen.getByPlaceholderText('Enter origin address')
    const destinationInput = screen.getByPlaceholderText('Enter destination address')
    const typeSelect = screen.getByRole('combobox')

    expect(dateInput).toBeDefined()
    await user.type(dateInput!, '2024-12-20')
    await user.type(originInput, 'Home')
    await user.type(destinationInput, 'Work')
    await user.selectOptions(typeSelect, 'single')

    // Submit the form
    const addTripButtons = screen.getAllByRole('button', { name: 'Add Trip' })
    // The first is the page button, the second is the form submit
    const submitButton = addTripButtons[1]
    await user.click(submitButton)

    // Wait for the API call to be made
    await waitFor(() => {
      expect(tripsApi.create).toHaveBeenCalledWith({
        date: '2024-12-20',
        origin: 'Home',
        destination: 'Work',
        type: 'single',
      })
    })
  })

  it('displays trips list when data is available', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('No trips recorded yet.')).toBeInTheDocument()
    })
  })

  it('shows empty state when no trips exist', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('No trips recorded yet.')).toBeInTheDocument()
      expect(screen.getByText('Add your first trip to get started!')).toBeInTheDocument()
    })
  })

  it('displays trips when data is available', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
      expect(screen.getAllByText('Work').length).toBeGreaterThan(0)
      expect(screen.getByText('Store')).toBeInTheDocument()
      expect(screen.getByText('Gym')).toBeInTheDocument()
    })
  })

  it('closes add trip form when cancel is clicked', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    // Wait for the component to load and show the Add Trip button
    await waitFor(() => {
      expect(screen.getByText('Add Trip')).toBeInTheDocument()
    })

    const addButton = screen.getByText('Add Trip')
    await user.click(addButton)

    const cancelButton = screen.getByText('Cancel')
    await user.click(cancelButton)

    expect(screen.queryByText('Add New Trip')).not.toBeInTheDocument()
  })

  describe('Search and Filtering', () => {
    it('renders search filter component', async () => {
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      await waitFor(() => {
        expect(screen.getByPlaceholderText('Search trips by origin, destination, or type...')).toBeInTheDocument()
        expect(screen.getByText('Advanced Filters')).toBeInTheDocument()
      })
    })

    it('filters trips by search term', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // Wait for trips to load
      await waitFor(() => {
        expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
        expect(screen.getAllByText('Work').length).toBeGreaterThan(0)
        expect(screen.getByText('Store')).toBeInTheDocument()
        expect(screen.getByText('Gym')).toBeInTheDocument()
      })

      const searchInput = screen.getByPlaceholderText('Search trips by origin, destination, or type...')
      await user.type(searchInput, 'home')

      // Wait for debounced search to complete
      await waitFor(() => {
        expect(screen.getByText('Showing 2 of 3 trips')).toBeInTheDocument()
      })

      // Verify filtered results - should show trips with "Home" but not "Store"
      expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
      expect(screen.getByText('Gym')).toBeInTheDocument()
      expect(screen.queryByText('Store')).not.toBeInTheDocument()
    })

    it('shows results summary when filters are active', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // Wait for trips to load
      await waitFor(() => {
        expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
      })

      const searchInput = screen.getByPlaceholderText('Search trips by origin, destination, or type...')
      await user.type(searchInput, 'home')

      await waitFor(() => {
        expect(screen.getByText('Showing 2 of 3 trips')).toBeInTheDocument()
      })
    })

    it('shows advanced filters when button is clicked', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      await waitFor(() => {
        expect(screen.getByText('Advanced Filters')).toBeInTheDocument()
      })

      const advancedButton = screen.getByText('Advanced Filters')
      await user.click(advancedButton)

      await waitFor(() => {
        expect(screen.getByText('Date From')).toBeInTheDocument()
        expect(screen.getByText('Date To')).toBeInTheDocument()
        expect(screen.getByText('Min Amount')).toBeInTheDocument()
        expect(screen.getByText('Max Amount')).toBeInTheDocument()
        expect(screen.getByText('Trip Type')).toBeInTheDocument()
      })
    })

    it('filters trips by trip type', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      await waitFor(() => {
        expect(screen.getByText('Advanced Filters')).toBeInTheDocument()
      })

      const advancedButton = screen.getByText('Advanced Filters')
      await user.click(advancedButton)

      const typeSelect = screen.getByLabelText('Trip Type')
      await user.selectOptions(typeSelect, 'single')

      await waitFor(() => {
        expect(screen.getByText('Showing 2 of 3 trips')).toBeInTheDocument()
      })

      // Verify filtered results - should show single trips but not round trips
      expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
      expect(screen.getByText('Gym')).toBeInTheDocument()
      expect(screen.queryByText('Store')).not.toBeInTheDocument()
    })

    it('shows appropriate message when no trips match filters', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // Wait for trips to load
      await waitFor(() => {
        expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
      })

      const searchInput = screen.getByPlaceholderText('Search trips by origin, destination, or type...')
      await user.type(searchInput, 'nonexistent')

      await waitFor(() => {
        expect(screen.getByText('No trips match your current filters.')).toBeInTheDocument()
        expect(screen.getByText('Try adjusting your search or filters.')).toBeInTheDocument()
      })
    })

    it('clears search when clear button is clicked', async () => {
      const user = userEvent.setup()
      const { tripsApi } = await import('../../lib/api')
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // Wait for trips to load
      await waitFor(() => {
        expect(screen.getAllByText('Home').length).toBeGreaterThan(0)
        expect(screen.getByText('Store')).toBeInTheDocument()
      })

      const searchInput = screen.getByPlaceholderText('Search trips by origin, destination, or type...')
      await user.type(searchInput, 'home')

      // Wait for search to filter results
      await waitFor(() => {
        expect(screen.getByText('Showing 2 of 3 trips')).toBeInTheDocument()
        expect(screen.queryByText('Store')).not.toBeInTheDocument()
      })

      // Find and click the clear button (X button in search input) - use a more specific selector
      const clearButtons = screen.getAllByRole('button')
      const clearButton = clearButtons.find(button => 
        button.closest('.relative') && button.querySelector('svg')
      )
      expect(clearButton).toBeDefined()
      await user.click(clearButton!)

      // Wait for the search to be cleared and all trips to show
      await waitFor(() => {
        expect(screen.queryByText('Showing 2 of 3 trips')).not.toBeInTheDocument()
        expect(screen.getByText('Store')).toBeInTheDocument()
      })
    })
  })
}) 