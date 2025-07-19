import React, { useState, useCallback, useMemo, useEffect, useRef } from 'react'
import { Search, Filter, X, Calendar, DollarSign, MapPin } from 'lucide-react'
import { debounce } from 'lodash'

export interface FilterOptions {
  search: string
  dateFrom: string
  dateTo: string
  minAmount: string
  maxAmount: string
  type: 'all' | 'single' | 'round'
  category: string
}

interface SearchFilterProps {
  onFilterChange: (filters: FilterOptions) => void
  filters: FilterOptions
  showAdvanced?: boolean
  placeholder?: string
}

const SearchFilter: React.FC<SearchFilterProps> = ({
  onFilterChange,
  filters,
  showAdvanced = false,
  placeholder = 'Search trips and expenses...'
}) => {
  const [isAdvancedOpen, setIsAdvancedOpen] = useState(showAdvanced)
  const [localSearch, setLocalSearch] = useState(filters.search)

  // Sync localSearch with filters.search when filters change externally
  useEffect(() => {
    setLocalSearch(filters.search)
  }, [filters.search])

  // Create a stable debounced function using useRef
  const debouncedSearchRef = useRef(
    debounce((searchTerm: string, currentFilters: FilterOptions) => {
      onFilterChange({ ...currentFilters, search: searchTerm })
    }, 300)
  )

  // Update the debounced function when onFilterChange changes
  useEffect(() => {
    debouncedSearchRef.current = debounce((searchTerm: string, currentFilters: FilterOptions) => {
      onFilterChange({ ...currentFilters, search: searchTerm })
    }, 300)
  }, [onFilterChange])

  const handleSearchChange = (value: string) => {
    setLocalSearch(value)
    debouncedSearchRef.current(value, filters)
  }

  const handleFilterChange = (key: keyof FilterOptions, value: string) => {
    onFilterChange({ ...filters, [key]: value })
  }

  const clearFilters = () => {
    // Cancel any pending debounced search calls
    if (debouncedSearchRef.current?.cancel) {
      debouncedSearchRef.current.cancel()
    }
    
    const clearedFilters: FilterOptions = {
      search: '',
      dateFrom: '',
      dateTo: '',
      minAmount: '',
      maxAmount: '',
      type: 'all',
      category: ''
    }
    setLocalSearch('')
    onFilterChange(clearedFilters)
  }

  const handleClearSearch = () => {
    // Cancel any pending debounced search calls
    if (debouncedSearchRef.current?.cancel) {
      debouncedSearchRef.current.cancel()
    }
    setLocalSearch('')
    onFilterChange({ ...filters, search: '' })
  }

  const hasActiveFilters = useMemo(() => {
    return (
      filters.search ||
      filters.dateFrom ||
      filters.dateTo ||
      filters.minAmount ||
      filters.maxAmount ||
      filters.type !== 'all' ||
      filters.category
    )
  }, [filters])

  return (
    <div className="space-y-4">
      {/* Search Bar */}
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          value={localSearch}
          onChange={(e) => handleSearchChange(e.target.value)}
          className="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
          placeholder={placeholder}
        />
        {localSearch && (
          <button
            onClick={handleClearSearch}
            className="absolute inset-y-0 right-0 pr-3 flex items-center"
          >
            <X className="h-5 w-5 text-gray-400 hover:text-gray-600" />
          </button>
        )}
      </div>

      {/* Filter Controls */}
      <div className="flex flex-wrap gap-3 items-center">
        <button
          onClick={() => setIsAdvancedOpen(!isAdvancedOpen)}
          className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
            isAdvancedOpen
              ? 'bg-primary-100 text-primary-700 border border-primary-200'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          <Filter className="h-4 w-4" />
          Advanced Filters
        </button>

        {hasActiveFilters && (
          <button
            onClick={clearFilters}
            className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-100 transition-colors"
          >
            <X className="h-4 w-4" />
            Clear All
          </button>
        )}
      </div>

      {/* Advanced Filters */}
      {isAdvancedOpen && (
        <div className="bg-gray-50 rounded-lg p-4 space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Date Range */}
            <div className="space-y-2">
              <label htmlFor="dateFrom" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <Calendar className="h-4 w-4" />
                Date From
              </label>
              <input
                id="dateFrom"
                type="date"
                value={filters.dateFrom}
                onChange={(e) => handleFilterChange('dateFrom', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="dateTo" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <Calendar className="h-4 w-4" />
                Date To
              </label>
              <input
                id="dateTo"
                type="date"
                value={filters.dateTo}
                onChange={(e) => handleFilterChange('dateTo', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              />
            </div>

            {/* Amount Range */}
            <div className="space-y-2">
              <label htmlFor="minAmount" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <DollarSign className="h-4 w-4" />
                Min Amount
              </label>
              <input
                id="minAmount"
                type="number"
                step="0.01"
                min="0"
                value={filters.minAmount}
                onChange={(e) => handleFilterChange('minAmount', e.target.value)}
                placeholder="0.00"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="maxAmount" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <DollarSign className="h-4 w-4" />
                Max Amount
              </label>
              <input
                id="maxAmount"
                type="number"
                step="0.01"
                min="0"
                value={filters.maxAmount}
                onChange={(e) => handleFilterChange('maxAmount', e.target.value)}
                placeholder="0.00"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {/* Trip Type */}
            <div className="space-y-2">
              <label htmlFor="tripType" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <MapPin className="h-4 w-4" />
                Trip Type
              </label>
              <select
                id="tripType"
                value={filters.type}
                onChange={(e) => handleFilterChange('type', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              >
                <option value="all">All Types</option>
                <option value="single">Single Trip</option>
                <option value="round">Round Trip</option>
              </select>
            </div>

            {/* Category */}
            <div className="space-y-2">
              <label htmlFor="category" className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <Filter className="h-4 w-4" />
                Category
              </label>
              <input
                id="category"
                type="text"
                value={filters.category}
                onChange={(e) => handleFilterChange('category', e.target.value)}
                placeholder="Enter category..."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm"
              />
            </div>
          </div>
        </div>
      )}

      {/* Active Filters Display */}
      {hasActiveFilters && (
        <div className="flex flex-wrap gap-2">
          {filters.search && (
            <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              Search: {filters.search}
              <button
                onClick={handleClearSearch}
                className="ml-1 hover:text-blue-600"
              >
                <X className="h-3 w-3" />
              </button>
            </span>
          )}
          {filters.dateFrom && (
            <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
              From: {filters.dateFrom}
              <button
                onClick={() => handleFilterChange('dateFrom', '')}
                className="ml-1 hover:text-green-600"
              >
                <X className="h-3 w-3" />
              </button>
            </span>
          )}
          {filters.dateTo && (
            <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
              To: {filters.dateTo}
              <button
                onClick={() => handleFilterChange('dateTo', '')}
                className="ml-1 hover:text-green-600"
              >
                <X className="h-3 w-3" />
              </button>
            </span>
          )}
          {filters.type !== 'all' && (
            <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
              Type: {filters.type}
              <button
                onClick={() => handleFilterChange('type', 'all')}
                className="ml-1 hover:text-purple-600"
              >
                <X className="h-3 w-3" />
              </button>
            </span>
          )}
        </div>
      )}
    </div>
  )
}

export default SearchFilter 