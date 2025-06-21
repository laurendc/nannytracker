import { rest } from 'msw'
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
  rest.get(`${baseUrl}/health`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        status: 'healthy',
        service: 'nannytracker-api',
      })
    )
  }),

  // Trips endpoints
  rest.get(`${baseUrl}/trips`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        trips: mockTrips,
        count: mockTrips.length,
      })
    )
  }),

  rest.post(`${baseUrl}/trips`, async (req, res, ctx) => {
    const newTrip = await req.json()
    const trip: Trip = {
      ...newTrip,
      miles: 5.0, // Mock calculated miles
    }
    return res(
      ctx.status(201),
      ctx.json(trip)
    )
  }),

  // Expenses endpoints
  rest.get(`${baseUrl}/expenses`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        expenses: mockExpenses,
        count: mockExpenses.length,
      })
    )
  }),

  rest.post(`${baseUrl}/expenses`, async (req, res, ctx) => {
    const expense = await req.json()
    return res(
      ctx.status(201),
      ctx.json(expense)
    )
  }),

  // Summaries endpoints
  rest.get(`${baseUrl}/summaries`, (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        summaries: mockSummaries,
        count: mockSummaries.length,
      })
    )
  }),
] 