import { http, HttpResponse } from 'msw'
import type { Trip, Expense, WeeklySummary } from '../../types'

const baseUrl = '/api'

// Mock data
const mockTrips: Trip[] = [
  {
    date: '2024-12-18',
    origin: 'Home',
    destination: 'Work',
    miles: 5.2,
    type: 'single',
  },
  {
    date: '2024-12-19',
    origin: 'Work',
    destination: 'Store',
    miles: 2.1,
    type: 'round',
  },
]

const mockExpenses: Expense[] = [
  {
    date: '2024-12-18',
    amount: 15.50,
    description: 'Lunch',
  },
  {
    date: '2024-12-19',
    amount: 8.75,
    description: 'Coffee',
  },
]

const mockSummaries: WeeklySummary[] = [
  {
    weekStart: '2024-12-15',
    weekEnd: '2024-12-21',
    totalMiles: 7.3,
    totalAmount: 5.11,
    trips: mockTrips,
    expenses: mockExpenses,
  },
]

export const handlers = [
  // Health check
  http.get(`${baseUrl}/health`, () => {
    return HttpResponse.json({
      status: 'healthy',
      service: 'nannytracker-api',
    }, { status: 200 })
  }),

  // Trips endpoints
  http.get(`${baseUrl}/trips`, () => {
    return HttpResponse.json({
      trips: mockTrips,
      count: mockTrips.length,
    }, { status: 200 })
  }),

  http.post(`${baseUrl}/trips`, async ({ request }) => {
    const newTrip = await request.json()
    const trip: Trip = {
      ...newTrip,
      miles: 5.0, // Mock calculated miles
    }
    return HttpResponse.json(trip, { status: 201 })
  }),

  // Expenses endpoints
  http.get(`${baseUrl}/expenses`, () => {
    return HttpResponse.json({
      expenses: mockExpenses,
      count: mockExpenses.length,
    }, { status: 200 })
  }),

  http.post(`${baseUrl}/expenses`, async ({ request }) => {
    const expense = await request.json()
    return HttpResponse.json(expense, { status: 201 })
  }),

  // Summaries endpoints
  http.get(`${baseUrl}/summaries`, () => {
    return HttpResponse.json({
      summaries: mockSummaries,
      count: mockSummaries.length,
    }, { status: 200 })
  }),
] 