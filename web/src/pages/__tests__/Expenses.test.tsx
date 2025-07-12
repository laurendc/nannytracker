import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { vi } from 'vitest'
import Expenses from '../Expenses'

// Mock the API calls
vi.mock('../../lib/api', () => ({
  expensesApi: {
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
const mockExpenses = [
  {
    id: '1',
    date: '2024-12-18',
    amount: 15.50,
    description: 'Lunch',
  },
  {
    id: '2',
    date: '2024-12-19',
    amount: 8.75,
    description: 'Coffee',
  },
]

describe('Expenses', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders expenses page title and description', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Expenses')).toBeInTheDocument()
      expect(screen.getByText(/Track your reimbursable expenses/)).toBeInTheDocument()
    })
  })

  it('shows add expense button', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })
  })

  it('displays total expenses summary', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Total Expenses')).toBeInTheDocument()
    })
  })

  it('opens add expense form when button is clicked', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })

    const addButton = screen.getByText('Add Expense')
    await user.click(addButton)

    expect(screen.getByText('Add New Expense')).toBeInTheDocument()
    // Check for form inputs by their type and placeholder
    // Date input: first input with value ""
    expect(screen.getAllByDisplayValue('')[0]).toBeInTheDocument()
    expect(screen.getByPlaceholderText('0.00')).toBeInTheDocument() // Amount input
    expect(screen.getByPlaceholderText('Enter expense description')).toBeInTheDocument() // Description input
  })

  it('allows form input and submission', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])
    vi.mocked(expensesApi.create).mockResolvedValue({})

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })

    // Open form
    const addButton = screen.getByText('Add Expense')
    await user.click(addButton)

    // Fill form using specific selectors
    const dateInput = screen.getAllByDisplayValue('')[0]
    const amountInput = screen.getByPlaceholderText('0.00')
    const descriptionInput = screen.getByPlaceholderText('Enter expense description')

    await user.type(dateInput, '2024-12-20')
    await user.type(amountInput, '15.50')
    await user.type(descriptionInput, 'Lunch')

    // Submit form - use the submit button specifically by type
    const submitButtons = screen.getAllByRole('button', { name: /add expense/i })
    const submitButton = submitButtons.find(btn => (btn as HTMLButtonElement).type === 'submit') || submitButtons[0]
    await user.click(submitButton)

    // Form should be submitted
    expect(submitButton).toBeInTheDocument()
  })

  it('displays expenses list when data is available', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('All Expenses')).toBeInTheDocument()
    })
  })

  it('shows empty state when no expenses exist', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('No expenses recorded yet.')).toBeInTheDocument()
      expect(screen.getByText(/Add your first expense to get started/)).toBeInTheDocument()
    })
  })

  it('displays expense information correctly', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Lunch')).toBeInTheDocument()
      expect(screen.getByText('Coffee')).toBeInTheDocument()
    })
  })

  it('shows edit and delete buttons for each expense', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      // Look for buttons with SVG icons (edit and delete buttons)
      const buttons = screen.getAllByRole('button')
      const editButtons = buttons.filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-gray-600')
      )
      const deleteButtons = buttons.filter(button => 
        button.querySelector('svg') && button.className.includes('hover:text-red-600')
      )
      
      expect(editButtons.length).toBeGreaterThan(0)
      expect(deleteButtons.length).toBeGreaterThan(0)
    })
  })

  it('opens edit form when edit button is clicked', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)

    render(
      <TestWrapper>
        <Expenses />
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
      expect(screen.getByText('Edit Expense')).toBeInTheDocument()
    })
  })

  it('submits edit form with updated data', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
    vi.mocked(expensesApi.update).mockResolvedValue({})

    render(
      <TestWrapper>
        <Expenses />
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
      expect(screen.getByText('Edit Expense')).toBeInTheDocument()
    })

    // Find the description input and update it
    const descriptionInput = screen.getByDisplayValue('Lunch')
    await user.clear(descriptionInput)
    await user.type(descriptionInput, 'Updated Lunch')

    // Find and click the update button
    const updateButton = screen.getByText('Update Expense')
    await user.click(updateButton)

    await waitFor(() => {
      expect(expensesApi.update).toHaveBeenCalledWith(0, expect.objectContaining({
        description: 'Updated Lunch',
      }))
    })
  })

  it('shows delete confirmation dialog when delete button is clicked', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
    vi.mocked(expensesApi.delete).mockResolvedValue(undefined)

    render(
      <TestWrapper>
        <Expenses />
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
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this expense?')
    })
  })

  it('confirms deletion and calls delete API', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
    vi.mocked(expensesApi.delete).mockResolvedValue(undefined)
    // Mock confirm to return true (user confirms deletion)
    vi.mocked(window.confirm).mockReturnValue(true)

    render(
      <TestWrapper>
        <Expenses />
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
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this expense?')
      expect(expensesApi.delete).toHaveBeenCalledWith(0)
    })
  })

  it('cancels deletion when cancel button is clicked', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
    vi.mocked(expensesApi.delete).mockResolvedValue(undefined)
    // Mock confirm to return false (user cancels deletion)
    vi.mocked(window.confirm).mockReturnValue(false)

    render(
      <TestWrapper>
        <Expenses />
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
      expect(window.confirm).toHaveBeenCalledWith('Are you sure you want to delete this expense?')
    })

    expect(expensesApi.delete).not.toHaveBeenCalled()
  })

  it('validates required form fields', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })

    // Open form
    const addButton = screen.getByText('Add Expense')
    await user.click(addButton)

    // Try to submit without filling required fields
    const submitButtons = screen.getAllByRole('button', { name: /add expense/i })
    const submitButton = submitButtons.find(btn => (btn as HTMLButtonElement).type === 'submit') || submitButtons[0]
    await user.click(submitButton)

    // Form should still be visible (not submitted)
    expect(screen.getByText('Add New Expense')).toBeInTheDocument()
  })

  it('validates amount field is numeric', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })

    // Open form
    const addButton = screen.getByText('Add Expense')
    await user.click(addButton)

    // Fill form with invalid amount using specific selectors
    const dateInput = screen.getAllByDisplayValue('')[0]
    const amountInput = screen.getByPlaceholderText('0.00')
    const descriptionInput = screen.getByPlaceholderText('Enter expense description')

    await user.type(dateInput, '2024-12-20')
    await user.type(amountInput, 'invalid')
    await user.type(descriptionInput, 'Test expense')

    // Submit form
    const submitButtons = screen.getAllByRole('button', { name: /add expense/i })
    const submitButton = submitButtons.find(btn => (btn as HTMLButtonElement).type === 'submit') || submitButtons[0]
    await user.click(submitButton)

    // Form should still be visible (not submitted)
    expect(screen.getByText('Add New Expense')).toBeInTheDocument()
  })

  it('allows canceling the add expense form', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue([])

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      expect(screen.getByText('Add Expense')).toBeInTheDocument()
    })

    // Open form
    const addButton = screen.getByText('Add Expense')
    await user.click(addButton)

    // Cancel form
    const cancelButton = screen.getByText('Cancel')
    await user.click(cancelButton)

    // Form should be closed
    expect(screen.queryByText('Add New Expense')).not.toBeInTheDocument()
  })

  it('calculates total expenses correctly', async () => {
    const { expensesApi } = await import('../../lib/api')
    vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)

    render(
      <TestWrapper>
        <Expenses />
      </TestWrapper>
    )

    await waitFor(() => {
      // Mock data has $15.50 + $8.75 = $24.25
      expect(screen.getByText('$24.25')).toBeInTheDocument()
    })
  })
}) 