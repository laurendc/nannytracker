import { useQuery } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Suspense, lazy } from 'react'
import { BarChart3, TrendingUp, DollarSign } from 'lucide-react'
import { summariesApi } from '../lib/api'
import LoadingSpinner from '../components/LoadingSpinner'

// Lazy load the Charts component to reduce initial bundle size
const Charts = lazy(() => import('../components/Charts'))

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6']

export default function Summaries() {
  const { data: summaries = [], isLoading } = useQuery({
    queryKey: ['summaries'],
    queryFn: summariesApi.getAll,
  })

  const chartData = summaries.map(summary => ({
    week: format(new Date(summary.weekStart), 'MMM d'),
    miles: summary.totalMiles,
    amount: summary.totalAmount,
    expenses: summary.expenses.reduce((sum, exp) => sum + exp.amount, 0),
  }))

  const totalMiles = summaries.reduce((sum, s) => sum + s.totalMiles, 0)
  const totalAmount = summaries.reduce((sum, s) => sum + s.totalAmount, 0)
  const totalExpenses = summaries.reduce((sum, s) => sum + s.expenses.reduce((eSum, e) => eSum + e.amount, 0), 0)

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="card">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded w-3/4"></div>
              </div>
            ))}
          </div>
          <div className="card">
            <div className="h-64 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Weekly Summaries</h1>
        <p className="text-gray-600 mt-1">
          Overview of your weekly mileage and expense tracking
        </p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <TrendingUp className="w-6 h-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Miles</p>
              <p className="text-2xl font-bold text-gray-900">{totalMiles.toFixed(1)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <DollarSign className="w-6 h-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Reimbursement</p>
              <p className="text-2xl font-bold text-gray-900">${totalAmount.toFixed(2)}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-yellow-100 rounded-lg">
              <BarChart3 className="w-6 h-6 text-yellow-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Expenses</p>
              <p className="text-2xl font-bold text-gray-900">${totalExpenses.toFixed(2)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Charts - Lazy loaded */}
      {summaries.length > 0 ? (
        <div className="space-y-6">
          <Suspense fallback={<LoadingSpinner size="lg" />}>
            <Charts chartData={chartData} />
          </Suspense>

          {/* Weekly Summaries List */}
          <div className="card">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Weekly Details</h2>
            <div className="space-y-4">
              {summaries.map((summary, index) => (
                <div key={index} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="font-semibold text-gray-900">
                      {format(new Date(summary.weekStart), 'MMM d')} - {format(new Date(summary.weekEnd), 'MMM d, yyyy')}
                    </h3>
                    <div className="text-right">
                      <p className="text-sm text-gray-600">Total Amount</p>
                      <p className="font-semibold text-gray-900">${summary.totalAmount.toFixed(2)}</p>
                    </div>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div>
                      <p className="text-sm text-gray-600">Total Miles</p>
                      <p className="font-medium text-gray-900">{summary.totalMiles.toFixed(1)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-600">Trips</p>
                      <p className="font-medium text-gray-900">{summary.trips.length}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-600">Expenses</p>
                      <p className="font-medium text-gray-900">{summary.expenses.length}</p>
                    </div>
                  </div>

                  {/* Trips in this week */}
                  {summary.trips.length > 0 && (
                    <div className="mb-3">
                      <p className="text-sm font-medium text-gray-700 mb-2">Trips:</p>
                      <div className="space-y-1">
                        {summary.trips.map((trip, tripIndex) => (
                          <div key={tripIndex} className="text-sm text-gray-600">
                            {format(new Date(trip.date), 'MMM d')}: {trip.origin} → {trip.destination} ({trip.miles} miles)
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Expenses in this week */}
                  {summary.expenses.length > 0 && (
                    <div>
                      <p className="text-sm font-medium text-gray-700 mb-2">Expenses:</p>
                      <div className="space-y-1">
                        {summary.expenses.map((expense, expenseIndex) => (
                          <div key={expenseIndex} className="text-sm text-gray-600">
                            {format(new Date(expense.date), 'MMM d')}: {expense.description} (${expense.amount.toFixed(2)})
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <div className="card">
          <div className="text-center py-12">
            <BarChart3 className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No weekly summaries available yet.</p>
            <p className="text-gray-400 text-sm mt-1">
              Add some trips and expenses to see your weekly summaries.
            </p>
          </div>
        </div>
      )}
    </div>
  )
} 