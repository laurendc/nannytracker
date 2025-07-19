import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Plus, Edit, Trash2, Receipt, Calendar, DollarSign } from 'lucide-react'
import { expensesApi } from '../lib/api'
import type { Expense } from '../types'
import SearchFilter, { type FilterOptions } from '../components/SearchFilter'
import { filterExpenses, getDefaultFilters, saveFiltersToLocalStorage, loadFiltersFromLocalStorage } from '../utils/filterUtils'

export default function Expenses() {
  const [isAddingExpense, setIsAddingExpense] = useState(false)
  const [editingExpense, setEditingExpense] = useState<{expense: Expense, index: number} | null>(null)
  const [filters, setFilters] = useState<FilterOptions>(() => loadFiltersFromLocalStorage('expenses-filters'))
  const [newExpense, setNewExpense] = useState<Partial<Expense>>({
    date: '',
    description: '',
    amount: 0,
  })

  const queryClient = useQueryClient()

  const { data: expenses = [], isLoading } = useQuery({
    queryKey: ['expenses'],
    queryFn: expensesApi.getAll,
  })

  // Filter expenses based on current filters
  const filteredExpenses = filterExpenses(expenses, filters)

  // Calculate total expenses
  const totalExpenses = expenses.reduce((sum, expense) => sum + expense.amount, 0)

  // Save filters to localStorage when they change
  useEffect(() => {
    saveFiltersToLocalStorage('expenses-filters', filters)
  }, [filters])

  const createExpenseMutation = useMutation({
    mutationFn: expensesApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      setIsAddingExpense(false)
      setNewExpense({ date: '', description: '', amount: 0 })
    },
  })

  const updateExpenseMutation = useMutation({
    mutationFn: ({ index, expense }: { index: number, expense: Expense }) => expensesApi.update(index, expense),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      setEditingExpense(null)
    },
  })

  const deleteExpenseMutation = useMutation({
    mutationFn: (index: number) => expensesApi.delete(index),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (newExpense.date && newExpense.description && newExpense.amount !== undefined) {
      createExpenseMutation.mutate({
        date: newExpense.date,
        description: newExpense.description,
        amount: newExpense.amount,
      })
    }
  }

  const handleEditSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingExpense) {
      updateExpenseMutation.mutate({
        index: editingExpense.index,
        expense: editingExpense.expense
      })
    }
  }

  const handleDelete = (index: number) => {
    if (confirm('Are you sure you want to delete this expense?')) {
      deleteExpenseMutation.mutate(index)
    }
  }

  const handleFilterChange = (newFilters: FilterOptions) => {
    setFilters(newFilters)
  }

  if (isLoading) {
    return (
      <div className="space-y-4 sm:space-y-6">
        <div className="animate-pulse">
          <div className="h-6 sm:h-8 bg-gray-200 rounded w-1/2 sm:w-1/4 mb-4 sm:mb-6"></div>
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="card">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* Mobile-first header */}
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Expenses</h1>
          <p className="text-sm sm:text-base text-gray-600 mt-1">
            Track your work-related expenses
          </p>
        </div>
        <button
          onClick={() => setIsAddingExpense(true)}
          className="btn btn-primary flex items-center justify-center w-full sm:w-auto touch-target"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Expense
        </button>
      </div>

      {/* Search and Filter */}
      <SearchFilter
        filters={filters}
        onFilterChange={handleFilterChange}
        placeholder="Search expenses by description or amount..."
      />

      {/* Total Expenses Summary */}
      <div className="card bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-gray-900">Total Expenses</h3>
            <p className="text-sm text-gray-600">All recorded expenses</p>
          </div>
          <div className="text-right">
            <p className="text-2xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
          </div>
        </div>
      </div>

      {/* Results Summary */}
      {filters.search || filters.dateFrom || filters.dateTo || filters.category ? (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600">
            Showing {filteredExpenses.length} of {expenses.length} expenses
          </p>
        </div>
      ) : null}

      {/* Add Expense Form - Mobile-optimized */}
      {isAddingExpense && (
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-4">Add New Expense</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="form-grid">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Date
                </label>
                <input
                  type="date"
                  value={newExpense.date}
                  onChange={(e) => setNewExpense({ ...newExpense, date: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Amount
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={newExpense.amount}
                  onChange={(e) => setNewExpense({ ...newExpense, amount: parseFloat(e.target.value) || 0 })}
                  className="input"
                  placeholder="0.00"
                  required
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <input
                type="text"
                value={newExpense.description}
                onChange={(e) => setNewExpense({ ...newExpense, description: e.target.value })}
                className="input"
                placeholder="Enter expense description"
                required
              />
            </div>
            <div className="flex flex-col sm:flex-row gap-3">
              <button type="submit" className="btn btn-primary touch-target" disabled={createExpenseMutation.isLoading}>
                {createExpenseMutation.isLoading ? 'Adding...' : 'Add Expense'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsAddingExpense(false)
                  setNewExpense({ date: '', description: '', amount: 0 })
                }}
                className="btn btn-secondary touch-target"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Edit Expense Form - Mobile-optimized */}
      {editingExpense && (
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-4">Edit Expense</h2>
          <form onSubmit={handleEditSubmit} className="space-y-4">
            <div className="form-grid">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Date
                </label>
                <input
                  type="date"
                  value={editingExpense.expense.date}
                  onChange={(e) => setEditingExpense({
                    ...editingExpense,
                    expense: { ...editingExpense.expense, date: e.target.value }
                  })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Amount
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={editingExpense.expense.amount}
                  onChange={(e) => setEditingExpense({
                    ...editingExpense,
                    expense: { ...editingExpense.expense, amount: parseFloat(e.target.value) || 0 }
                  })}
                  className="input"
                  placeholder="0.00"
                  required
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <input
                type="text"
                value={editingExpense.expense.description}
                onChange={(e) => setEditingExpense({
                  ...editingExpense,
                  expense: { ...editingExpense.expense, description: e.target.value }
                })}
                className="input"
                placeholder="Enter expense description"
                required
              />
            </div>
            <div className="flex flex-col sm:flex-row gap-3">
              <button type="submit" className="btn btn-primary touch-target" disabled={updateExpenseMutation.isLoading}>
                {updateExpenseMutation.isLoading ? 'Updating...' : 'Update Expense'}
              </button>
              <button
                type="button"
                onClick={() => setEditingExpense(null)}
                className="btn btn-secondary touch-target"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Expenses List - Mobile-first cards */}
      <div className="space-y-4">
        {filteredExpenses.length === 0 ? (
          <div className="card text-center py-8">
            <Receipt className="w-12 h-12 mx-auto text-gray-400 mb-4" />
            <p className="text-gray-500 text-sm sm:text-base">
              {expenses.length === 0 ? 'No expenses recorded yet.' : 'No expenses match your current filters.'}
            </p>
            <p className="text-gray-400 text-xs sm:text-sm mt-2">
              {expenses.length === 0 ? 'Add your first expense to get started!' : 'Try adjusting your search or filters.'}
            </p>
          </div>
        ) : (
          filteredExpenses.map((expense, index) => (
            <div key={index} className="card hover:shadow-md transition-shadow">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div className="flex-1 mb-4 sm:mb-0">
                  <div className="flex items-center mb-2">
                    <Calendar className="w-4 h-4 text-gray-400 mr-2" />
                    <span className="text-sm text-gray-600">
                      {format(new Date(expense.date), 'MMM d, yyyy')}
                    </span>
                  </div>
                  <div className="flex items-center mb-2">
                    <Receipt className="w-4 h-4 text-gray-400 mr-2" />
                    <span className="text-sm sm:text-base font-medium text-gray-900">
                      {expense.description}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <DollarSign className="w-4 h-4 text-gray-400 mr-2" />
                    <span className="text-sm sm:text-base font-semibold text-gray-900">
                      ${expense.amount.toFixed(2)}
                    </span>
                  </div>
                </div>
                <div className="flex flex-row sm:flex-col gap-2 sm:gap-3">
                  <button
                    onClick={() => setEditingExpense({ expense, index })}
                    className="btn btn-secondary text-sm px-3 py-2 flex items-center justify-center touch-target"
                  >
                    <Edit className="w-4 h-4 mr-1" />
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(index)}
                    className="btn bg-red-100 text-red-700 hover:bg-red-200 text-sm px-3 py-2 flex items-center justify-center touch-target"
                    disabled={deleteExpenseMutation.isLoading}
                  >
                    <Trash2 className="w-4 h-4 mr-1" />
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
} 