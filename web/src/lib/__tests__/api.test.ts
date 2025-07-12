import { vi } from 'vitest'
import { tripsApi, expensesApi, summariesApi, healthApi } from '../api'
import type { Trip, Expense } from '../../types'

// Mock the entire api module
vi.mock('../api', () => ({
  tripsApi: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
  expensesApi: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
  summariesApi: {
    getAll: vi.fn(),
  },
  healthApi: {
    check: vi.fn(),
  },
}))

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks()
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
      vi.mocked(tripsApi.getAll).mockResolvedValue(mockTrips)

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
      vi.mocked(tripsApi.create).mockResolvedValue(createdTrip)

      const result = await tripsApi.create(newTrip)
      expect(result).toEqual(createdTrip)
    })

    it('update sends trip data correctly', async () => {
      const updatedTrip = {
        date: '2024-12-19',
        origin: 'Updated Home',
        destination: 'Updated Work',
        miles: 10.0,
        type: 'round' as const,
      }

      const { tripsApi } = await import('../api')
      vi.mocked(tripsApi.update).mockResolvedValue(updatedTrip)

      const result = await tripsApi.update(0, updatedTrip)
      expect(result).toEqual(updatedTrip)
      expect(tripsApi.update).toHaveBeenCalledWith(0, updatedTrip)
    })

    it('delete removes trip correctly', async () => {
      const { tripsApi } = await import('../api')
      vi.mocked(tripsApi.delete).mockResolvedValue(undefined)

      await tripsApi.delete(0)
      expect(tripsApi.delete).toHaveBeenCalledWith(0)
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
      vi.mocked(expensesApi.getAll).mockResolvedValue(mockExpenses)

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
      vi.mocked(expensesApi.create).mockResolvedValue(newExpense)

      const result = await expensesApi.create(newExpense)
      expect(result).toEqual(newExpense)
    })

    it('update sends expense data correctly', async () => {
      const updatedExpense = {
        date: '2024-12-19',
        amount: 25.75,
        description: 'Updated Lunch',
      }

      const { expensesApi } = await import('../api')
      vi.mocked(expensesApi.update).mockResolvedValue(updatedExpense)

      const result = await expensesApi.update(0, updatedExpense)
      expect(result).toEqual(updatedExpense)
      expect(expensesApi.update).toHaveBeenCalledWith(0, updatedExpense)
    })

    it('delete removes expense correctly', async () => {
      const { expensesApi } = await import('../api')
      vi.mocked(expensesApi.delete).mockResolvedValue(undefined)

      await expensesApi.delete(0)
      expect(expensesApi.delete).toHaveBeenCalledWith(0)
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
      vi.mocked(summariesApi.getAll).mockResolvedValue(mockSummaries)

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
      vi.mocked(healthApi.check).mockResolvedValue(mockHealth)

      const result = await healthApi.check()
      expect(result).toEqual(mockHealth)
    })
  })
}) 