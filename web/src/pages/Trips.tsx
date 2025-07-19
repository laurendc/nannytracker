import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Plus, Edit, Trash2, Car, MapPin, Calendar, ArrowRight } from 'lucide-react'
import { tripsApi } from '../lib/api'
import type { Trip } from '../types'
import SearchFilter, { type FilterOptions } from '../components/SearchFilter'
import { filterTrips, getDefaultFilters, saveFiltersToLocalStorage, loadFiltersFromLocalStorage } from '../utils/filterUtils'

export default function Trips() {
  const [isAddingTrip, setIsAddingTrip] = useState(false)
  const [editingTrip, setEditingTrip] = useState<{trip: Trip, index: number} | null>(null)
  const [filters, setFilters] = useState<FilterOptions>(() => loadFiltersFromLocalStorage('trips-filters'))
  const [newTrip, setNewTrip] = useState<Partial<Trip>>({
    date: '',
    origin: '',
    destination: '',
    type: 'single',
  })

  const queryClient = useQueryClient()

  const { data: trips = [], isLoading } = useQuery({
    queryKey: ['trips'],
    queryFn: tripsApi.getAll,
  })

  // Filter trips based on current filters
  const filteredTrips = filterTrips(trips, filters)

  // Save filters to localStorage when they change
  useEffect(() => {
    saveFiltersToLocalStorage('trips-filters', filters)
  }, [filters])

  const createTripMutation = useMutation({
    mutationFn: tripsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] })
      setIsAddingTrip(false)
      setNewTrip({ date: '', origin: '', destination: '', type: 'single' })
    },
  })

  const updateTripMutation = useMutation({
    mutationFn: ({ index, trip }: { index: number, trip: Trip }) => tripsApi.update(index, trip),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] })
      setEditingTrip(null)
    },
  })

  const deleteTripMutation = useMutation({
    mutationFn: (index: number) => tripsApi.delete(index),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] })
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (newTrip.date && newTrip.origin && newTrip.destination && newTrip.type) {
      createTripMutation.mutate({
        date: newTrip.date,
        origin: newTrip.origin,
        destination: newTrip.destination,
        type: newTrip.type,
      })
    }
  }

  const handleEditSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingTrip) {
      updateTripMutation.mutate({
        index: editingTrip.index,
        trip: editingTrip.trip
      })
    }
  }

  const handleDelete = (index: number) => {
    if (confirm('Are you sure you want to delete this trip?')) {
      deleteTripMutation.mutate(index)
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
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Trips</h1>
          <p className="text-sm sm:text-base text-gray-600 mt-1">
            Manage your mileage tracking entries
          </p>
        </div>
        <button
          onClick={() => setIsAddingTrip(true)}
          className="btn btn-primary flex items-center justify-center w-full sm:w-auto touch-target"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Trip
        </button>
      </div>

      {/* Search and Filter */}
      <SearchFilter
        filters={filters}
        onFilterChange={handleFilterChange}
        placeholder="Search trips by origin, destination, or type..."
      />

      {/* Results Summary */}
      {filters.search || filters.dateFrom || filters.dateTo || filters.type !== 'all' ? (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600">
            Showing {filteredTrips.length} of {trips.length} trips
          </p>
        </div>
      ) : null}

      {/* Add Trip Form - Mobile-optimized */}
      {isAddingTrip && (
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-4">Add New Trip</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="form-grid">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Date
                </label>
                <input
                  type="date"
                  value={newTrip.date}
                  onChange={(e) => setNewTrip({ ...newTrip, date: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Type
                </label>
                <select
                  value={newTrip.type}
                  onChange={(e) => setNewTrip({ ...newTrip, type: e.target.value as 'single' | 'round' })}
                  className="input"
                  required
                >
                  <option value="single">Single Trip</option>
                  <option value="round">Round Trip</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Origin
              </label>
              <input
                type="text"
                value={newTrip.origin}
                onChange={(e) => setNewTrip({ ...newTrip, origin: e.target.value })}
                className="input"
                placeholder="Enter origin address"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Destination
              </label>
              <input
                type="text"
                value={newTrip.destination}
                onChange={(e) => setNewTrip({ ...newTrip, destination: e.target.value })}
                className="input"
                placeholder="Enter destination address"
                required
              />
            </div>
            <div className="flex flex-col sm:flex-row gap-3">
              <button type="submit" className="btn btn-primary touch-target" disabled={createTripMutation.isLoading}>
                {createTripMutation.isLoading ? 'Adding...' : 'Add Trip'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsAddingTrip(false)
                  setNewTrip({ date: '', origin: '', destination: '', type: 'single' })
                }}
                className="btn btn-secondary touch-target"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Edit Trip Form - Mobile-optimized */}
      {editingTrip && (
        <div className="card">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900 mb-4">Edit Trip</h2>
          <form onSubmit={handleEditSubmit} className="space-y-4">
            <div className="form-grid">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Date
                </label>
                <input
                  type="date"
                  value={editingTrip.trip.date}
                  onChange={(e) => setEditingTrip({
                    ...editingTrip,
                    trip: { ...editingTrip.trip, date: e.target.value }
                  })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Type
                </label>
                <select
                  value={editingTrip.trip.type}
                  onChange={(e) => setEditingTrip({
                    ...editingTrip,
                    trip: { ...editingTrip.trip, type: e.target.value as 'single' | 'round' }
                  })}
                  className="input"
                  required
                >
                  <option value="single">Single Trip</option>
                  <option value="round">Round Trip</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Origin
              </label>
              <input
                type="text"
                value={editingTrip.trip.origin}
                onChange={(e) => setEditingTrip({
                  ...editingTrip,
                  trip: { ...editingTrip.trip, origin: e.target.value }
                })}
                className="input"
                placeholder="Enter origin address"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Destination
              </label>
              <input
                type="text"
                value={editingTrip.trip.destination}
                onChange={(e) => setEditingTrip({
                  ...editingTrip,
                  trip: { ...editingTrip.trip, destination: e.target.value }
                })}
                className="input"
                placeholder="Enter destination address"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Miles
              </label>
              <input
                type="number"
                step="0.1"
                value={editingTrip.trip.miles}
                onChange={(e) => setEditingTrip({
                  ...editingTrip,
                  trip: { ...editingTrip.trip, miles: parseFloat(e.target.value) || 0 }
                })}
                className="input"
                placeholder="Enter miles"
                required
              />
            </div>
            <div className="flex flex-col sm:flex-row gap-3">
              <button type="submit" className="btn btn-primary touch-target" disabled={updateTripMutation.isLoading}>
                {updateTripMutation.isLoading ? 'Updating...' : 'Update Trip'}
              </button>
              <button
                type="button"
                onClick={() => setEditingTrip(null)}
                className="btn btn-secondary touch-target"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Trips List - Mobile-first cards */}
      <div className="space-y-4">
        {filteredTrips.length === 0 ? (
          <div className="card text-center py-8">
            <Car className="w-12 h-12 mx-auto text-gray-400 mb-4" />
            <p className="text-gray-500 text-sm sm:text-base">
              {trips.length === 0 ? 'No trips recorded yet.' : 'No trips match your current filters.'}
            </p>
            <p className="text-gray-400 text-xs sm:text-sm mt-2">
              {trips.length === 0 ? 'Add your first trip to get started!' : 'Try adjusting your search or filters.'}
            </p>
          </div>
        ) : (
          filteredTrips.map((trip, index) => (
            <div key={index} className="card hover:shadow-md transition-shadow">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div className="flex-1 mb-4 sm:mb-0">
                  <div className="flex items-center mb-2">
                    <Calendar className="w-4 h-4 text-gray-400 mr-2" />
                    <span className="text-sm text-gray-600">
                      {format(new Date(trip.date), 'MMM d, yyyy')}
                    </span>
                    <span className="ml-2 inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {trip.type}
                    </span>
                  </div>
                  <div className="flex items-center mb-2">
                    <MapPin className="w-4 h-4 text-gray-400 mr-2" />
                    <div className="flex items-center text-sm sm:text-base" data-testid="trip-row">
                      <span className="font-medium text-gray-900 truncate">{trip.origin}</span>
                      <ArrowRight className="w-4 h-4 mx-2 text-gray-400" />
                      <span className="font-medium text-gray-900 truncate">{trip.destination}</span>
                    </div>
                  </div>
                  <div className="flex items-center">
                    <Car className="w-4 h-4 text-gray-400 mr-2" />
                    <span className="text-sm sm:text-base font-semibold text-gray-900">
                      {trip.miles} miles
                    </span>
                  </div>
                </div>
                <div className="flex flex-row sm:flex-col gap-2 sm:gap-3">
                  <button
                    onClick={() => setEditingTrip({ trip, index })}
                    className="btn btn-secondary text-sm px-3 py-2 flex items-center justify-center touch-target"
                  >
                    <Edit className="w-4 h-4 mr-1" />
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(index)}
                    className="btn bg-red-100 text-red-700 hover:bg-red-200 text-sm px-3 py-2 flex items-center justify-center touch-target"
                    disabled={deleteTripMutation.isLoading}
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