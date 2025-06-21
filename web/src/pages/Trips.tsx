import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { format } from 'date-fns'
import { Plus, Edit, Trash2, Car } from 'lucide-react'
import { tripsApi } from '../lib/api'
import type { Trip } from '../types'

export default function Trips() {
  const [isAddingTrip, setIsAddingTrip] = useState(false)
  const [editingTrip, setEditingTrip] = useState<Trip | null>(null)
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

  const createTripMutation = useMutation({
    mutationFn: tripsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trips'] })
      setIsAddingTrip(false)
      setNewTrip({ date: '', origin: '', destination: '', type: 'single' })
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
          <h1 className="text-2xl font-bold text-gray-900">Trips</h1>
          <p className="text-gray-600 mt-1">
            Manage your mileage tracking entries
          </p>
        </div>
        <button
          onClick={() => setIsAddingTrip(true)}
          className="btn btn-primary flex items-center"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Trip
        </button>
      </div>

      {/* Add Trip Form */}
      {isAddingTrip && (
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Add New Trip</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
                <label className="block text-sm font-medium text-gray-700 mb-1">
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
              <label className="block text-sm font-medium text-gray-700 mb-1">
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
              <label className="block text-sm font-medium text-gray-700 mb-1">
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
            <div className="flex gap-3">
              <button type="submit" className="btn btn-primary">
                {createTripMutation.isLoading ? 'Adding...' : 'Add Trip'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setIsAddingTrip(false)
                  setNewTrip({ date: '', origin: '', destination: '', type: 'single' })
                }}
                className="btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Trips List */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">All Trips</h2>
        {trips.length === 0 ? (
          <div className="text-center py-12">
            <Car className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">No trips recorded yet.</p>
            <p className="text-gray-400 text-sm mt-1">
              Add your first trip to get started.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {trips.map((trip, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center space-x-4">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <Car className="w-5 h-5 text-blue-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">
                      {trip.origin} → {trip.destination}
                    </p>
                    <p className="text-sm text-gray-600">
                      {format(new Date(trip.date), 'MMM d, yyyy')} • {trip.type} • {trip.miles} miles
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => setEditingTrip(trip)}
                    className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
                  >
                    <Edit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => {
                      if (confirm('Are you sure you want to delete this trip?')) {
                        // TODO: Implement delete functionality
                      }
                    }}
                    className="p-2 text-gray-400 hover:text-red-600 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
} 