import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Plus, Edit, Trash2, Receipt } from 'lucide-react'
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
    if (newExpense.date && newExpense.amount && newExpense.description) {
      createExpenseMutation.mutate({
        date: newExpense.date,
        amount: newExpense.amount,
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
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
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
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Expenses</h1>
          <p className="text-gray-600 mt-1">
            Track your reimbursable expenses
          </p>
        </div>
        <button
          onClick={() => setIsAddingExpense(true)}
          className="btn btn-primary flex items-center"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Expense
        </button>
      </div>

      {/* Total Expenses Summary */}
      <div className="card">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-600">Total Expenses</p>
            <p className="text-3xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
          </div>
          <div className="p-3 bg-yellow-100 rounded-lg">
            <Receipt className="w-8 h-8 text-yellow-600" />
          </div>
        </div>
      </div>

      {/* Add Expense Form */}
      {isAddingExpense && (
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Add New Expense</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
              <label className="block text-sm font-medium text-gray-700 mb-1">
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
            <div className="flex gap-3">
              <button type="submit" className="btn btn-primary" disabled={createExpenseMutation.isLoading}>
                {createExpenseMutation.isLoading ? 'Adding...' : 'Add Expense'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsAddingExpense(false)
                  setNewExpense({ date: '', amount: 0, description: '' })
                }}
                className="btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Edit Expense Form */}
      {editingExpense && (
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Edit Expense</h2>
          <form onSubmit={handleEditSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
              <label className="block text-sm font-medium text-gray-700 mb-1">
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
            <div className="flex gap-3">
              <button type="submit" className="btn btn-primary" disabled={updateExpenseMutation.isLoading}>
                {updateExpenseMutation.isLoading ? 'Updating...' : 'Update Expense'}
              </button>
              <button
                type="button"
                onClick={() => setEditingExpense(null)}
                className="btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Expenses List */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">All Expenses</h2>
        {expenses.length === 0 ? (
          <div className="text-center py-12">
            <Receipt className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No expenses recorded yet.</p>
            <p className="text-gray-400 text-sm mt-1">
              Add your first expense to get started.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {expenses.map((expense, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center space-x-4">
                  <div className="p-2 bg-yellow-100 rounded-lg">
                    <Receipt className="w-5 h-5 text-yellow-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{expense.description}</p>
                    <p className="text-sm text-gray-600">
                      {format(new Date(expense.date), 'MMM d, yyyy')}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  <div className="text-right">
                    <p className="font-semibold text-gray-900">${expense.amount.toFixed(2)}</p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => setEditingExpense({ expense, index })}
                      className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
                      disabled={updateExpenseMutation.isLoading}
                    >
                      <Edit className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(index)}
                      className="p-2 text-gray-400 hover:text-red-600 transition-colors"
                      disabled={deleteExpenseMutation.isLoading}
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
} 