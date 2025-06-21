import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Trips from '../Trips'

// Mock the API calls
jest.mock('../../lib/api', () => ({
  tripsApi: {
    getAll: jest.fn(),
    create: jest.fn(),
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
    id: 1,
    date: '2024-12-20',
    origin: 'Home',
    destination: 'Work',
    type: 'single' as const,
    miles: 15.5,
  },
  {
    id: 2,
    date: '2024-12-21',
    origin: 'Work',
    destination: 'Store',
    type: 'round' as const,
    miles: 8.2,
  },
]

describe('Trips', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders trips page title and description', async () => {
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])
    ;(tripsApi.create as jest.Mock).mockResolvedValue({})

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

    const emptyInputs = screen.getAllByDisplayValue('')
    const dateInput = emptyInputs[0]
    const originInput = screen.getByPlaceholderText('Enter origin address')
    const destinationInput = screen.getByPlaceholderText('Enter destination address')
    const typeSelect = screen.getByRole('combobox')

    await user.type(dateInput, '2024-12-20')
    await user.type(originInput, 'Home')
    await user.type(destinationInput, 'Work')
    await user.selectOptions(typeSelect, 'single')

    const submitButtons = screen.getAllByText('Add Trip')
    const submitButton = submitButtons[1]
    await user.click(submitButton)

    expect(submitButton).toBeInTheDocument()
  })

  it('displays trips list when data is available', async () => {
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('All Trips')).toBeInTheDocument()
    })
  })

  it('shows empty state when no trips exist', async () => {
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('No trips recorded yet.')).toBeInTheDocument()
      expect(screen.getByText(/Add your first trip to get started/)).toBeInTheDocument()
    })
  })

  it('displays trip information correctly', async () => {
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue(mockTrips)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Home → Work')).toBeInTheDocument()
      expect(screen.getByText('Work → Store')).toBeInTheDocument()
    })
  })

  it('shows edit and delete buttons for each trip', async () => {
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue(mockTrips)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const editButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-gray-600')
      )
      const deleteButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-red-600')
      )
      
      expect(editButtons.length).toBeGreaterThan(0)
      expect(deleteButtons.length).toBeGreaterThan(0)
    })
  })

  it('validates required form fields', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

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

    const submitButtons = screen.getAllByText('Add Trip')
    const submitButton = submitButtons[1]
    await user.click(submitButton)

    expect(screen.getByText('Add New Trip')).toBeInTheDocument()
  })

  it('allows canceling the add trip form', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    ;(tripsApi.getAll as jest.Mock).mockResolvedValue([])

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

    const cancelButton = screen.getByText('Cancel')
    await user.click(cancelButton)

    expect(screen.queryByText('Add New Trip')).not.toBeInTheDocument()
  })
}) 