import { describe, it, expect, beforeEach } from 'vitest'
import { filterTrips, filterExpenses, getDefaultFilters, saveFiltersToLocalStorage, loadFiltersFromLocalStorage } from '../filterUtils'
import type { Trip, Expense } from '../../types'
import type { FilterOptions } from '../../components/SearchFilter'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('filterUtils', () => {
  const mockTrips: Trip[] = [
    {
      date: '2024-01-15',
      origin: 'Home',
      destination: 'Work',
      type: 'single',
      miles: 15.5
    },
    {
      date: '2024-01-20',
      origin: 'Work',
      destination: 'Store',
      type: 'round',
      miles: 8.2
    },
    {
      date: '2024-02-01',
      origin: 'Home',
      destination: 'Gym',
      type: 'single',
      miles: 5.0
    }
  ]

  const mockExpenses: Expense[] = [
    {
      date: '2024-01-15',
      description: 'Lunch',
      amount: 12.50
    },
    {
      date: '2024-01-20',
      description: 'Coffee',
      amount: 3.75
    },
    {
      date: '2024-02-01',
      description: 'Gas',
      amount: 45.00
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('filterTrips', () => {
    it('returns all trips when no filters are applied', () => {
      const filters: FilterOptions = getDefaultFilters()
      const result = filterTrips(mockTrips, filters)
      expect(result).toEqual(mockTrips)
    })

    it('filters trips by search term in origin', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'home' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].origin).toBe('Home')
      expect(result[1].origin).toBe('Home')
    })

    it('filters trips by search term in destination', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'work' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].destination).toBe('Work')
      expect(result[1].origin).toBe('Work')
    })

    it('filters trips by search term in type', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'single' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].type).toBe('single')
      expect(result[1].type).toBe('single')
    })

    it('filters trips by date range - from date', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), dateFrom: '2024-01-20' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].date).toBe('2024-01-20')
      expect(result[1].date).toBe('2024-02-01')
    })

    it('filters trips by date range - to date', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), dateTo: '2024-01-20' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].date).toBe('2024-01-15')
      expect(result[1].date).toBe('2024-01-20')
    })

    it('filters trips by date range - both from and to', () => {
      const filters: FilterOptions = { 
        ...getDefaultFilters(), 
        dateFrom: '2024-01-15', 
        dateTo: '2024-01-20' 
      }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].date).toBe('2024-01-15')
      expect(result[1].date).toBe('2024-01-20')
    })

    it('filters trips by type', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), type: 'single' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].type).toBe('single')
      expect(result[1].type).toBe('single')
    })

    it('filters trips by minimum miles', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), minAmount: '10' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(1)
      expect(result[0].miles).toBe(15.5)
    })

    it('filters trips by maximum miles', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), maxAmount: '10' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].miles).toBe(8.2)
      expect(result[1].miles).toBe(5.0)
    })

    it('combines multiple filters', () => {
      const filters: FilterOptions = { 
        ...getDefaultFilters(), 
        search: 'home',
        type: 'single',
        minAmount: '10'
      }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(1)
      expect(result[0].origin).toBe('Home')
      expect(result[0].type).toBe('single')
      expect(result[0].miles).toBe(15.5)
    })

    it('returns empty array when no trips match filters', () => {
      const filters: FilterOptions = { 
        ...getDefaultFilters(), 
        search: 'nonexistent'
      }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(0)
    })
  })

  describe('filterExpenses', () => {
    it('returns all expenses when no filters are applied', () => {
      const filters: FilterOptions = getDefaultFilters()
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toEqual(mockExpenses)
    })

    it('filters expenses by search term in description', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'lunch' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].description).toBe('Lunch')
    })

    it('filters expenses by search term in amount', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: '12' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].amount).toBe(12.50)
    })

    it('filters expenses by date range - from date', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), dateFrom: '2024-01-20' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(2)
      expect(result[0].date).toBe('2024-01-20')
      expect(result[1].date).toBe('2024-02-01')
    })

    it('filters expenses by date range - to date', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), dateTo: '2024-01-20' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(2)
      expect(result[0].date).toBe('2024-01-15')
      expect(result[1].date).toBe('2024-01-20')
    })

    it('filters expenses by minimum amount', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), minAmount: '10' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(2)
      expect(result[0].amount).toBe(12.50)
      expect(result[1].amount).toBe(45.00)
    })

    it('filters expenses by maximum amount', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), maxAmount: '10' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].amount).toBe(3.75)
    })

    it('filters expenses by category (description contains)', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), category: 'coffee' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].description).toBe('Coffee')
    })

    it('combines multiple filters for expenses', () => {
      const filters: FilterOptions = { 
        ...getDefaultFilters(), 
        search: 'lunch',
        minAmount: '10'
      }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].description).toBe('Lunch')
      expect(result[0].amount).toBe(12.50)
    })

    it('returns empty array when no expenses match filters', () => {
      const filters: FilterOptions = { 
        ...getDefaultFilters(), 
        search: 'nonexistent'
      }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(0)
    })
  })

  describe('getDefaultFilters', () => {
    it('returns default filter values', () => {
      const result = getDefaultFilters()
      expect(result).toEqual({
        search: '',
        dateFrom: '',
        dateTo: '',
        minAmount: '',
        maxAmount: '',
        type: 'all',
        category: ''
      })
    })
  })

  describe('localStorage functions', () => {
    it('saves filters to localStorage', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'test' }
      saveFiltersToLocalStorage('test-key', filters)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('test-key', JSON.stringify(filters))
    })

    it('loads filters from localStorage', () => {
      const savedFilters = { search: 'test', type: 'single' }
      localStorageMock.getItem.mockReturnValue(JSON.stringify(savedFilters))
      
      const result = loadFiltersFromLocalStorage('test-key')
      expect(localStorageMock.getItem).toHaveBeenCalledWith('test-key')
      expect(result).toEqual({ ...getDefaultFilters(), ...savedFilters })
    })

    it('returns default filters when localStorage is empty', () => {
      localStorageMock.getItem.mockReturnValue(null)
      
      const result = loadFiltersFromLocalStorage('test-key')
      expect(result).toEqual(getDefaultFilters())
    })

    it('returns default filters when localStorage has invalid JSON', () => {
      localStorageMock.getItem.mockReturnValue('invalid json')
      
      const result = loadFiltersFromLocalStorage('test-key')
      expect(result).toEqual(getDefaultFilters())
    })

    it('handles localStorage errors gracefully', () => {
      localStorageMock.setItem.mockImplementation(() => {
        throw new Error('Storage quota exceeded')
      })
      
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'test' }
      expect(() => saveFiltersToLocalStorage('test-key', filters)).not.toThrow()
    })
  })

  describe('edge cases', () => {
    it('handles empty arrays', () => {
      const filters: FilterOptions = getDefaultFilters()
      expect(filterTrips([], filters)).toEqual([])
      expect(filterExpenses([], filters)).toEqual([])
    })

    it('handles case-insensitive search', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'HOME' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(2)
      expect(result[0].origin).toBe('Home')
    })

    it('handles partial matches in search', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), search: 'cof' }
      const result = filterExpenses(mockExpenses, filters)
      expect(result).toHaveLength(1)
      expect(result[0].description).toBe('Coffee')
    })

    it('handles decimal amounts in filters', () => {
      const filters: FilterOptions = { ...getDefaultFilters(), minAmount: '10.5' }
      const result = filterTrips(mockTrips, filters)
      expect(result).toHaveLength(1)
      expect(result[0].miles).toBe(15.5)
    })
  })
}) 