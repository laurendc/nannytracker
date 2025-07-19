import type { Trip, Expense } from '../types'
import type { FilterOptions } from '../components/SearchFilter'

export function filterTrips(trips: Trip[], filters: FilterOptions): Trip[] {
  return trips.filter(trip => {
    // Search filter
    if (filters.search) {
      const searchTerm = filters.search.toLowerCase()
      const matchesSearch = 
        trip.origin.toLowerCase().includes(searchTerm) ||
        trip.destination.toLowerCase().includes(searchTerm) ||
        trip.type.toLowerCase().includes(searchTerm)
      
      if (!matchesSearch) return false
    }

    // Date range filter
    if (filters.dateFrom && trip.date < filters.dateFrom) {
      return false
    }
    if (filters.dateTo && trip.date > filters.dateTo) {
      return false
    }

    // Trip type filter
    if (filters.type !== 'all' && trip.type !== filters.type) {
      return false
    }

    // Amount range filter (for miles)
    if (filters.minAmount && trip.miles < parseFloat(filters.minAmount)) {
      return false
    }
    if (filters.maxAmount && trip.miles > parseFloat(filters.maxAmount)) {
      return false
    }

    return true
  })
}

export function filterExpenses(expenses: Expense[], filters: FilterOptions): Expense[] {
  return expenses.filter(expense => {
    // Search filter
    if (filters.search) {
      const searchTerm = filters.search.toLowerCase()
      const matchesSearch = 
        expense.description.toLowerCase().includes(searchTerm) ||
        expense.amount.toString().includes(searchTerm)
      
      if (!matchesSearch) return false
    }

    // Date range filter
    if (filters.dateFrom && expense.date < filters.dateFrom) {
      return false
    }
    if (filters.dateTo && expense.date > filters.dateTo) {
      return false
    }

    // Amount range filter
    if (filters.minAmount && expense.amount < parseFloat(filters.minAmount)) {
      return false
    }
    if (filters.maxAmount && expense.amount > parseFloat(filters.maxAmount)) {
      return false
    }

    // Category filter
    if (filters.category) {
      const categoryTerm = filters.category.toLowerCase()
      const matchesCategory = expense.description.toLowerCase().includes(categoryTerm)
      if (!matchesCategory) return false
    }

    return true
  })
}

export function getDefaultFilters(): FilterOptions {
  return {
    search: '',
    dateFrom: '',
    dateTo: '',
    minAmount: '',
    maxAmount: '',
    type: 'all',
    category: ''
  }
}

export function saveFiltersToLocalStorage(key: string, filters: FilterOptions): void {
  try {
    localStorage.setItem(key, JSON.stringify(filters))
  } catch (error) {
    console.warn('Failed to save filters to localStorage:', error)
  }
}

export function loadFiltersFromLocalStorage(key: string): FilterOptions {
  try {
    const saved = localStorage.getItem(key)
    if (saved) {
      const parsed = JSON.parse(saved)
      return { ...getDefaultFilters(), ...parsed }
    }
  } catch (error) {
    console.warn('Failed to load filters from localStorage:', error)
  }
  return getDefaultFilters()
} 