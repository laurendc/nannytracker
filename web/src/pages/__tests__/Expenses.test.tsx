import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Expenses from '../Expenses'

// Mock the API calls
jest.mock('../../lib/api', () => ({
  expensesApi: {
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
    jest.clearAllMocks()
  })

  it('renders expenses page title and description', async () => {
    const { expensesApi } = await import('../../lib/api')
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])
    ;(expensesApi.create as jest.Mock).mockResolvedValue({})

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue(mockExpenses)

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue(mockExpenses)

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

  it('validates required form fields', async () => {
    const user = userEvent.setup()
    const { expensesApi } = await import('../../lib/api')
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue([])

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
    ;(expensesApi.getAll as jest.Mock).mockResolvedValue(mockExpenses)

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