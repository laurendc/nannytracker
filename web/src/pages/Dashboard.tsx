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
      <div className="space-y-4 sm:space-y-6">
        <div className="animate-pulse">
          <div className="h-6 sm:h-8 bg-gray-200 rounded w-1/2 sm:w-1/4 mb-4 sm:mb-6"></div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="card">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-6 sm:h-8 bg-gray-200 rounded w-3/4"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      <div>
        <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm sm:text-base text-gray-600 mt-1">
          Overview of your mileage and expense tracking
        </p>
      </div>

      {/* Stats Cards - Mobile-first grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
        <div className="card">
          <div className="flex items-center">
            <div className="p-2 sm:p-3 bg-blue-100 rounded-lg touch-target">
              <Car className="w-5 h-5 sm:w-6 sm:h-6 text-blue-600" />
            </div>
            <div className="ml-3 sm:ml-4">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Total Trips</p>
              <p className="text-xl sm:text-2xl font-bold text-gray-900">{trips.length}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 sm:p-3 bg-green-100 rounded-lg touch-target">
              <TrendingUp className="w-5 h-5 sm:w-6 sm:h-6 text-green-600" />
            </div>
            <div className="ml-3 sm:ml-4">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Total Miles</p>
              <p className="text-xl sm:text-2xl font-bold text-gray-900">{totalMiles.toFixed(1)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 sm:p-3 bg-yellow-100 rounded-lg touch-target">
              <Receipt className="w-5 h-5 sm:w-6 sm:h-6 text-yellow-600" />
            </div>
            <div className="ml-3 sm:ml-4">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Total Expenses</p>
              <p className="text-xl sm:text-2xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 sm:p-3 bg-purple-100 rounded-lg touch-target">
              <DollarSign className="w-5 h-5 sm:w-6 sm:h-6 text-purple-600" />
            </div>
            <div className="ml-3 sm:ml-4">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Reimbursement</p>
              <p className="text-xl sm:text-2xl font-bold text-gray-900">
                ${(totalMiles * 0.7).toFixed(2)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Activity - Mobile-first layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
        {/* Recent Trips */}
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Recent Trips</h2>
          {trips.length === 0 ? (
            <p className="text-sm sm:text-base text-gray-500 text-center py-8">No trips recorded yet.</p>
          ) : (
            <div className="space-y-3">
              {trips.slice(-5).reverse().map((trip, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-gray-900 text-sm sm:text-base truncate">
                      {trip.origin} → {trip.destination}
                    </p>
                    <p className="text-xs sm:text-sm text-gray-600 mt-1">
                      {trip.date && !isNaN(new Date(trip.date).getTime())
                        ? format(new Date(trip.date), 'MMM d, yyyy')
                        : 'Invalid date'} • {trip.type}
                    </p>
                  </div>
                  <div className="text-right ml-4 flex-shrink-0">
                    <p className="font-medium text-gray-900 text-sm sm:text-base">{trip.miles} miles</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Recent Expenses */}
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Recent Expenses</h2>
          {expenses.length === 0 ? (
            <p className="text-sm sm:text-base text-gray-500 text-center py-8">No expenses recorded yet.</p>
          ) : (
            <div className="space-y-3">
              {expenses.slice(-5).reverse().map((expense, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-gray-900 text-sm sm:text-base truncate">{expense.description}</p>
                    <p className="text-xs sm:text-sm text-gray-600 mt-1">
                      {expense.date && !isNaN(new Date(expense.date).getTime())
                        ? format(new Date(expense.date), 'MMM d, yyyy')
                        : 'Invalid date'}
                    </p>
                  </div>
                  <div className="text-right ml-4 flex-shrink-0">
                    <p className="font-medium text-gray-900 text-sm sm:text-base">${expense.amount.toFixed(2)}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Current Week Summary - Mobile-optimized */}
      {currentWeekSummary && (
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Current Week Summary</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6">
            <div className="text-center sm:text-left">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Week</p>
              <p className="text-sm sm:text-lg font-semibold text-gray-900 mt-1">
                {currentWeekSummary && currentWeekSummary.weekStart && !isNaN(new Date(currentWeekSummary.weekStart).getTime())
                  ? format(new Date(currentWeekSummary.weekStart), 'MMM d')
                  : 'Invalid date'}
                {' - '}
                {currentWeekSummary && currentWeekSummary.weekEnd && !isNaN(new Date(currentWeekSummary.weekEnd).getTime())
                  ? format(new Date(currentWeekSummary.weekEnd), 'MMM d, yyyy')
                  : 'Invalid date'}
              </p>
            </div>
            <div className="text-center sm:text-left">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Total Miles</p>
              <p className="text-sm sm:text-lg font-semibold text-gray-900 mt-1">
                {typeof currentWeekSummary.totalMiles === 'number'
                  ? currentWeekSummary.totalMiles.toFixed(1)
                  : 'N/A'}
              </p>
            </div>
            <div className="text-center sm:text-left">
              <p className="text-xs sm:text-sm font-medium text-gray-600">Total Amount</p>
              <p className="text-sm sm:text-lg font-semibold text-gray-900 mt-1">
                {typeof currentWeekSummary.totalAmount === 'number'
                  ? `$${currentWeekSummary.totalAmount.toFixed(2)}`
                  : 'N/A'}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
} 