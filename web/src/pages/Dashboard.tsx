import { useQuery } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Car, Receipt, DollarSign, TrendingUp } from 'lucide-react'
import { tripsApi, expensesApi, summariesApi } from '../lib/api'

export default function Dashboard() {
  const { data: trips = [], isLoading: tripsLoading } = useQuery({
    queryKey: ['trips'],
    queryFn: tripsApi.getAll,
  })

  const { data: expenses = [], isLoading: expensesLoading } = useQuery({
    queryKey: ['expenses'],
    queryFn: expensesApi.getAll,
  })

  const { data: summaries = [], isLoading: summariesLoading } = useQuery({
    queryKey: ['summaries'],
    queryFn: summariesApi.getAll,
  })

  const isLoading = tripsLoading || expensesLoading || summariesLoading

  const totalMiles = trips.reduce((sum, trip) => sum + trip.miles, 0)
  const totalExpenses = expenses.reduce((sum, expense) => sum + expense.amount, 0)
  const currentWeekSummary = summaries[summaries.length - 1]

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="card">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded w-3/4"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-1">
          Overview of your mileage and expense tracking
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <Car className="w-6 h-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Trips</p>
              <p className="text-2xl font-bold text-gray-900">{trips.length}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <TrendingUp className="w-6 h-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Miles</p>
              <p className="text-2xl font-bold text-gray-900">{totalMiles.toFixed(1)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-yellow-100 rounded-lg">
              <Receipt className="w-6 h-6 text-yellow-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Expenses</p>
              <p className="text-2xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-purple-100 rounded-lg">
              <DollarSign className="w-6 h-6 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Reimbursement</p>
              <p className="text-2xl font-bold text-gray-900">
                ${(totalMiles * 0.7).toFixed(2)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Trips */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Trips</h2>
          {trips.length === 0 ? (
            <p className="text-gray-500">No trips recorded yet.</p>
          ) : (
            <div className="space-y-3">
              {trips.slice(-5).reverse().map((trip, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="font-medium text-gray-900">
                      {trip.origin} → {trip.destination}
                    </p>
                    <p className="text-sm text-gray-600">
                      {format(new Date(trip.date), 'MMM d, yyyy')} • {trip.type}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-gray-900">{trip.miles} miles</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Recent Expenses */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Expenses</h2>
          {expenses.length === 0 ? (
            <p className="text-gray-500">No expenses recorded yet.</p>
          ) : (
            <div className="space-y-3">
              {expenses.slice(-5).reverse().map((expense, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="font-medium text-gray-900">{expense.description}</p>
                    <p className="text-sm text-gray-600">
                      {format(new Date(expense.date), 'MMM d, yyyy')}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-gray-900">${expense.amount.toFixed(2)}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Current Week Summary */}
      {currentWeekSummary && (
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Current Week Summary</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <p className="text-sm font-medium text-gray-600">Week</p>
              <p className="text-lg font-semibold text-gray-900">
                {format(new Date(currentWeekSummary.weekStart), 'MMM d')} - {format(new Date(currentWeekSummary.weekEnd), 'MMM d, yyyy')}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-600">Total Miles</p>
              <p className="text-lg font-semibold text-gray-900">{currentWeekSummary.totalMiles.toFixed(1)}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-600">Total Amount</p>
              <p className="text-lg font-semibold text-gray-900">${currentWeekSummary.totalAmount.toFixed(2)}</p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
} 