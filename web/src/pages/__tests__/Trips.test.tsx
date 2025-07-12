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
    vi.clearAllMocks()
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
    await act(async () => {
      await user.click(addButton)
    })

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

    // After successful submission, the form should close and the API should be called
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
      expect(screen.getByText('All Trips')).toBeInTheDocument()
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
      expect(screen.getByText(/Add your first trip to get started/)).toBeInTheDocument()
    })
  })

  it('displays trip information correctly', async () => {
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

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
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

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

  it('opens edit form when edit button is clicked', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const editButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-gray-600')
      )
      expect(editButtons.length).toBeGreaterThan(0)
    })

    const editButton = screen.getAllByRole('button').filter(button => 
      button.querySelector('svg') && button.className.includes('hover:text-gray-600')
    )[0]
    
    await user.click(editButton)

    await waitFor(() => {
      expect(screen.getByText('Edit Trip')).toBeInTheDocument()
    })
  })

  it('submits edit form with updated data', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
    vi.mocked(tripsApi.update).mockResolvedValue({
      date: '2024-12-18',
      origin: 'Updated Home',
      destination: 'Work',
      type: 'single',
      miles: 15.5
    })

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const editButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-gray-600')
      )
      expect(editButtons.length).toBeGreaterThan(0)
    })

    const editButton = screen.getAllByRole('button').filter(button => 
      button.querySelector('svg') && button.className.includes('hover:text-gray-600')
    )[0]
    
    await user.click(editButton)

    await waitFor(() => {
      expect(screen.getByText('Edit Trip')).toBeInTheDocument()
    })

    // Find the origin input and update it
    const originInput = screen.getByDisplayValue('Home')
    await user.clear(originInput)
    await user.type(originInput, 'Updated Home')

    // Find and click the update button
    const updateButton = screen.getByText('Update Trip')
    await user.click(updateButton)

    await waitFor(() => {
      expect(tripsApi.update).toHaveBeenCalledWith(0, expect.objectContaining({
        origin: 'Updated Home',
      }))
    })
  })

  it('shows delete confirmation dialog when delete button is clicked', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
    vi.mocked(tripsApi.delete).mockResolvedValue(undefined)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const deleteButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-red-600')
      )
      expect(deleteButtons.length).toBeGreaterThan(0)
    })

    const deleteButton = screen.getAllByRole('button').filter(button => 
      button.querySelector('svg') && button.className.includes('hover:text-red-600')
    )[0]
    
    await user.click(deleteButton)

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this trip?')
    })
  })

  it('confirms deletion and calls delete API', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
    vi.mocked(tripsApi.delete).mockResolvedValue(undefined)
    // Mock confirm to return true (user confirms deletion)
    vi.mocked(window.confirm).mockReturnValue(true)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const deleteButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-red-600')
      )
      expect(deleteButtons.length).toBeGreaterThan(0)
    })

    const deleteButton = screen.getAllByRole('button').filter(button => 
      button.querySelector('svg') && button.className.includes('hover:text-red-600')
    )[0]
    
    await user.click(deleteButton)

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this trip?')
      expect(tripsApi.delete).toHaveBeenCalledWith(0)
    })
  })

  it('cancels deletion when cancel button is clicked', async () => {
    const user = userEvent.setup()
    const { tripsApi } = await import('../../lib/api')
    vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
    vi.mocked(tripsApi.delete).mockResolvedValue(undefined)
    // Mock confirm to return false (user cancels deletion)
    vi.mocked(window.confirm).mockReturnValue(false)

    render(
      <TestWrapper>
        <Trips />
      </TestWrapper>
    )

    await waitFor(() => {
      const deleteButtons = screen.getAllByRole('button').filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-red-600')
      )
      expect(deleteButtons.length).toBeGreaterThan(0)
    })

    const deleteButton = screen.getAllByRole('button').filter(button => 
      button.querySelector('svg') && button.className.includes('hover:text-red-600')
    )[0]
    
    await user.click(deleteButton)

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this trip?')
    })

    expect(tripsApi.delete).not.toHaveBeenCalled()
  })

  it('validates required form fields', async () => {
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
    await user.click(addButton)

    const submitButtons = screen.getAllByText('Add Trip')
    const submitButton = submitButtons[1]
    await user.click(submitButton)

    expect(screen.getByText('Add New Trip')).toBeInTheDocument()
  })

  it('allows canceling the add trip form', async () => {
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
    await user.click(addButton)

    const cancelButton = screen.getByText('Cancel')
    await user.click(cancelButton)

    expect(screen.queryByText('Add New Trip')).not.toBeInTheDocument()
  })
}) 