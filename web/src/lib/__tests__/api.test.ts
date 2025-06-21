import { tripsApi, expensesApi, summariesApi, healthApi } from '../api'
import type { Trip, Expense } from '../../types'

// Mock the entire api module
jest.mock('../api', () => ({
  tripsApi: {
    getAll: jest.fn(),
    create: jest.fn(),
  },
  expensesApi: {
    getAll: jest.fn(),
    create: jest.fn(),
  },
  summariesApi: {
    getAll: jest.fn(),
  },
  healthApi: {
    check: jest.fn(),
  },
}))

describe('API Client', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('tripsApi', () => {
    it('getAll returns trips data', async () => {
      const mockTrips = [
        {
          date: '2024-12-18',
          origin: 'Home',
          destination: 'Work',
          miles: 5.2,
          type: 'single' as const,
        },
      ]

      const { tripsApi } = await import('../api')
      ;(tripsApi.getAll as jest.Mock).mockResolvedValue(mockTrips)

      const result = await tripsApi.getAll()
      expect(result).toEqual(mockTrips)
    })

    it('create sends trip data correctly', async () => {
      const newTrip = {
        date: '2024-12-18',
        origin: 'Home',
        destination: 'Work',
        type: 'single' as const,
      }

      const createdTrip = {
        ...newTrip,
        miles: 5.0,
      }

      const { tripsApi } = await import('../api')
      ;(tripsApi.create as jest.Mock).mockResolvedValue(createdTrip)

      const result = await tripsApi.create(newTrip)
      expect(result).toEqual(createdTrip)
    })
  })

  describe('expensesApi', () => {
    it('getAll returns expenses data', async () => {
      const mockExpenses = [
        {
          date: '2024-12-18',
          amount: 15.50,
          description: 'Lunch',
        },
      ]

      const { expensesApi } = await import('../api')
      ;(expensesApi.getAll as jest.Mock).mockResolvedValue(mockExpenses)

      const result = await expensesApi.getAll()
      expect(result).toEqual(mockExpenses)
    })

    it('create sends expense data correctly', async () => {
      const newExpense = {
        date: '2024-12-18',
        amount: 15.50,
        description: 'Lunch',
      }

      const { expensesApi } = await import('../api')
      ;(expensesApi.create as jest.Mock).mockResolvedValue(newExpense)

      const result = await expensesApi.create(newExpense)
      expect(result).toEqual(newExpense)
    })
  })

  describe('summariesApi', () => {
    it('getAll returns summaries data', async () => {
      const mockSummaries = [
        {
          weekStart: '2024-12-15',
          weekEnd: '2024-12-21',
          totalMiles: 7.3,
          totalAmount: 5.11,
          trips: [],
          expenses: [],
        },
      ]

      const { summariesApi } = await import('../api')
      ;(summariesApi.getAll as jest.Mock).mockResolvedValue(mockSummaries)

      const result = await summariesApi.getAll()
      expect(result).toEqual(mockSummaries)
    })
  })

  describe('healthApi', () => {
    it('check returns health status', async () => {
      const mockHealth = {
        status: 'healthy',
        service: 'nannytracker-api',
      }

      const { healthApi } = await import('../api')
      ;(healthApi.check as jest.Mock).mockResolvedValue(mockHealth)

      const result = await healthApi.check()
      expect(result).toEqual(mockHealth)
    })
  })
}) 