import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { vi } from 'vitest'
import Trips from '../Trips'
import Expenses from '../Expenses'
import { tripsApi, expensesApi } from '../../lib/api'

// Mock the API calls
vi.mock('../../lib/api', () => ({
  tripsApi: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
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

describe('CRUD Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Trip CRUD Operations', () => {
    it('completes full CRUD cycle for trips', async () => {
      const user = userEvent.setup()
      const mockTrips = [
        {
          id: 1,
          date: '2024-12-20',
          origin: 'Home',
          destination: 'Work',
          type: 'single' as const,
          miles: 15.5,
        },
      ]

      // Mock API responses
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)
      vi.mocked(tripsApi.create).mockResolvedValue({
        id: 2,
        date: '2024-12-21',
        origin: 'Home',
        destination: 'Store',
        type: 'round' as const,
        miles: 10.0,
      })
      vi.mocked(tripsApi.update).mockResolvedValue({
        id: 1,
        date: '2024-12-20',
        origin: 'Updated Home',
        destination: 'Updated Work',
        type: 'single' as const,
        miles: 20.0,
      })
      vi.mocked(tripsApi.delete).mockResolvedValue(undefined)

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // READ - verify initial data loads
      await waitFor(() => {
        expect(screen.getByText('Home → Work')).toBeInTheDocument()
      })

      // CREATE - add new trip
      const addButton = screen.getByText('Add Trip')
      await user.click(addButton)

      await waitFor(() => {
        expect(screen.getByText('Add New Trip')).toBeInTheDocument()
      })

      // Fill form
      const emptyInputs = screen.getAllByDisplayValue('')
      const dateInput = emptyInputs[0]
      const originInput = screen.getByPlaceholderText('Enter origin address')
      const destinationInput = screen.getByPlaceholderText('Enter destination address')
      const typeSelect = screen.getByRole('combobox')

      await user.type(dateInput, '2024-12-21')
      await user.type(originInput, 'Home')
      await user.type(destinationInput, 'Store')
      await user.selectOptions(typeSelect, 'round')

      // Submit form
      const submitButtons = screen.getAllByText('Add Trip')
      const submitButton = submitButtons[1]
      await user.click(submitButton)

      // Verify API was called
      expect(tripsApi.create).toHaveBeenCalledWith({
        date: '2024-12-21',
        origin: 'Home',
        destination: 'Store',
        type: 'round',
      })

      // UPDATE - edit existing trip
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

      // Update origin
      const originUpdateInput = screen.getByDisplayValue('Home')
      await user.clear(originUpdateInput)
      await user.type(originUpdateInput, 'Updated Home')

      // Submit update
      const updateButton = screen.getByText('Update Trip')
      await user.click(updateButton)

      // Verify update API was called
      expect(tripsApi.update).toHaveBeenCalledWith(0, expect.objectContaining({
        origin: 'Updated Home',
      }))

      // DELETE - remove trip
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
        expect(screen.getByText('Are you sure you want to delete this trip?')).toBeInTheDocument()
      })

      // Confirm deletion
      const confirmButton = screen.getByText('Delete')
      await user.click(confirmButton)

      // Verify delete API was called
      expect(tripsApi.delete).toHaveBeenCalledWith(0)
    })

    it('handles error scenarios gracefully', async () => {
      const user = userEvent.setup()
      
      // Mock API failures
      vi.mocked(tripsApi.getAll).mockRejectedValue(new Error('Network error'))
      vi.mocked(tripsApi.create).mockRejectedValue(new Error('Creation failed'))
      vi.mocked(tripsApi.update).mockRejectedValue(new Error('Update failed'))
      vi.mocked(tripsApi.delete).mockRejectedValue(new Error('Delete failed'))

      render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      // Should handle initial load failure gracefully
      await waitFor(() => {
        expect(screen.getByText('All Trips')).toBeInTheDocument()
      })

      // Test creation error handling
      const addButton = screen.getByText('Add Trip')
      await user.click(addButton)

      const emptyInputs = screen.getAllByDisplayValue('')
      const dateInput = emptyInputs[0]
      const originInput = screen.getByPlaceholderText('Enter origin address')
      const destinationInput = screen.getByPlaceholderText('Enter destination address')
      const typeSelect = screen.getByRole('combobox')

      await user.type(dateInput, '2024-12-21')
      await user.type(originInput, 'Test')
      await user.type(destinationInput, 'Test')
      await user.selectOptions(typeSelect, 'single')

      const submitButtons = screen.getAllByText('Add Trip')
      const submitButton = submitButtons[1]
      await user.click(submitButton)

      // Should handle creation error gracefully
      expect(tripsApi.create).toHaveBeenCalled()
    })
  })

  describe('Expense CRUD Operations', () => {
    it('completes full CRUD cycle for expenses', async () => {
      const user = userEvent.setup()
      const mockExpenses = [
        {
          id: '1',
          date: '2024-12-18',
          amount: 15.50,
          description: 'Lunch',
        },
      ]

      // Mock API responses
      vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)
      vi.mocked(expensesApi.create).mockResolvedValue({
        id: '2',
        date: '2024-12-19',
        amount: 8.75,
        description: 'Coffee',
      })
      vi.mocked(expensesApi.update).mockResolvedValue({
        id: '1',
        date: '2024-12-18',
        amount: 25.00,
        description: 'Updated Lunch',
      })
      vi.mocked(expensesApi.delete).mockResolvedValue(undefined)

      render(
        <TestWrapper>
          <Expenses />
        </TestWrapper>
      )

      // READ - verify initial data loads
      await waitFor(() => {
        expect(screen.getByText('Lunch')).toBeInTheDocument()
      })

      // CREATE - add new expense
      const addButton = screen.getByText('Add Expense')
      await user.click(addButton)

      await waitFor(() => {
        expect(screen.getByText('Add New Expense')).toBeInTheDocument()
      })

      // Fill form
      const dateInput = screen.getAllByDisplayValue('')[0]
      const amountInput = screen.getByPlaceholderText('0.00')
      const descriptionInput = screen.getByPlaceholderText('Enter expense description')

      await user.type(dateInput, '2024-12-19')
      await user.type(amountInput, '8.75')
      await user.type(descriptionInput, 'Coffee')

      // Submit form
      const submitButtons = screen.getAllByRole('button', { name: /add expense/i })
      const submitButton = submitButtons.find(btn => (btn as HTMLButtonElement).type === 'submit') || submitButtons[0]
      await user.click(submitButton)

      // Verify API was called
      expect(expensesApi.create).toHaveBeenCalledWith({
        date: '2024-12-19',
        amount: 8.75,
        description: 'Coffee',
      })

      // UPDATE - edit existing expense
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

      // Update description
      const descriptionUpdateInput = screen.getByDisplayValue('Lunch')
      await user.clear(descriptionUpdateInput)
      await user.type(descriptionUpdateInput, 'Updated Lunch')

      // Submit update
      const updateButton = screen.getByText('Update Expense')
      await user.click(updateButton)

      // Verify update API was called
      expect(expensesApi.update).toHaveBeenCalledWith(0, expect.objectContaining({
        description: 'Updated Lunch',
      }))

      // DELETE - remove expense
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
        expect(screen.getByText('Are you sure you want to delete this expense?')).toBeInTheDocument()
      })

      // Confirm deletion
      const confirmButton = screen.getByText('Delete')
      await user.click(confirmButton)

      // Verify delete API was called
      expect(expensesApi.delete).toHaveBeenCalledWith(0)
    })

    it('validates form data before submission', async () => {
      const user = userEvent.setup()
      
      vi.mocked(expensesApi.getAll).mockResolvedValue([])
      vi.mocked(expensesApi.create).mockResolvedValue({})

      render(
        <TestWrapper>
          <Expenses />
        </TestWrapper>
      )

      const addButton = screen.getByText('Add Expense')
      await user.click(addButton)

      // Test empty form submission
      const submitButtons = screen.getAllByRole('button', { name: /add expense/i })
      const submitButton = submitButtons.find(btn => (btn as HTMLButtonElement).type === 'submit') || submitButtons[0]
      await user.click(submitButton)

      // Form should still be visible (validation failed)
      expect(screen.getByText('Add New Expense')).toBeInTheDocument()

      // Test invalid amount
      const dateInput = screen.getAllByDisplayValue('')[0]
      const amountInput = screen.getByPlaceholderText('0.00')
      const descriptionInput = screen.getByPlaceholderText('Enter expense description')

      await user.type(dateInput, '2024-12-19')
      await user.type(amountInput, 'invalid')
      await user.type(descriptionInput, 'Test')

      await user.click(submitButton)

      // Form should still be visible (validation failed)
      expect(screen.getByText('Add New Expense')).toBeInTheDocument()

      // API should not be called with invalid data
      expect(expensesApi.create).not.toHaveBeenCalled()
    })
  })

  describe('Cross-Entity Operations', () => {
    it('maintains data consistency across different entity types', async () => {
      const user = userEvent.setup()
      
      // Mock both trips and expenses data
      vi.mocked(tripsApi.getAll).mockResolvedValue([
        {
          id: 1,
          date: '2024-12-20',
          origin: 'Home',
          destination: 'Work',
          type: 'single' as const,
          miles: 15.5,
        },
      ])
      
      vi.mocked(expensesApi.getAll).mockResolvedValue([
        {
          id: '1',
          date: '2024-12-20',
          amount: 15.50,
          description: 'Lunch',
        },
      ])

      // Test trips component
      const { rerender } = render(
        <TestWrapper>
          <Trips />
        </TestWrapper>
      )

      await waitFor(() => {
        expect(screen.getByText('Home → Work')).toBeInTheDocument()
      })

      // Switch to expenses component
      rerender(
        <TestWrapper>
          <Expenses />
        </TestWrapper>
      )

      await waitFor(() => {
        expect(screen.getByText('Lunch')).toBeInTheDocument()
      })

      // Verify both API calls were made
      expect(tripsApi.getAll).toHaveBeenCalled()
      expect(expensesApi.getAll).toHaveBeenCalled()
    })
  })
}) 