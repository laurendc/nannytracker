import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Plus, Edit, Trash2, Receipt, Calendar, DollarSign } from 'lucide-react'
import { expensesApi } from '../lib/api'
import type { Expense } from '../types'

export default function Expenses() {
  const [isAddingExpense, setIsAddingExpense] = useState(false)
  const [editingExpense, setEditingExpense] = useState<{expense: Expense, index: number} | null>(null)
  const [newExpense, setNewExpense] = useState<Partial<Expense>>({
    date: '',
    amount: 0,
    description: '',
  })

  const queryClient = useQueryClient()

  const { data: expenses = [], isLoading } = useQuery({
    queryKey: ['expenses'],
    queryFn: expensesApi.getAll,
  })

  const createExpenseMutation = useMutation({
    mutationFn: expensesApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['expenses'] })
      setIsAddingExpense(false)
      setNewExpense({ date: '', amount: 0, description: '' })
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
    const amount = newExpense.amount || 0
    if (newExpense.date && amount > 0 && newExpense.description) {
      createExpenseMutation.mutate({
        date: newExpense.date,
        amount: amount,
        description: newExpense.description,
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

  const totalExpenses = expenses.reduce((sum, expense) => sum + expense.amount, 0)

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
            Track your reimbursable expenses
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

      {/* Total Expenses Summary - Mobile-optimized */}
      <div className="card">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-xs sm:text-sm font-medium text-gray-600">Total Expenses</p>
            <p className="text-2xl sm:text-3xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
          </div>
          <div className="p-2 sm:p-3 bg-yellow-100 rounded-lg touch-target">
            <Receipt className="w-6 h-6 sm:w-8 sm:h-8 text-yellow-600" />
          </div>
        </div>
      </div>

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
                  value={newExpense.amount === 0 ? '' : newExpense.amount}
                  onChange={(e) => {
                    const value = e.target.value
                    const numValue = value === '' ? 0 : parseFloat(value)
                    setNewExpense({ ...newExpense, amount: isNaN(numValue) ? 0 : numValue })
                  }}
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
                  setNewExpense({ date: '', amount: 0, description: '' })
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
        {expenses.length === 0 ? (
          <div className="card text-center py-8">
            <Receipt className="w-12 h-12 mx-auto text-gray-400 mb-4" />
            <p className="text-gray-500 text-sm sm:text-base">No expenses recorded yet.</p>
            <p className="text-gray-400 text-xs sm:text-sm mt-2">Add your first expense to get started!</p>
          </div>
        ) : (
          expenses.map((expense, index) => (
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
                    <span className="font-medium text-gray-900 text-sm sm:text-base truncate">
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