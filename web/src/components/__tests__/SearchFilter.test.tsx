import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SearchFilter, { type FilterOptions } from '../SearchFilter'

// Mock lodash debounce
vi.mock('lodash', () => ({
  debounce: (fn: Function) => fn
}))

describe('SearchFilter', () => {
  const mockOnFilterChange = vi.fn()
  const defaultFilters: FilterOptions = {
    search: '',
    dateFrom: '',
    dateTo: '',
    minAmount: '',
    maxAmount: '',
    type: 'all',
    category: ''
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders search input with placeholder', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        placeholder="Test placeholder"
      />
    )

    expect(screen.getByPlaceholderText('Test placeholder')).toBeInTheDocument()
  })

  it('renders with default placeholder when none provided', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByPlaceholderText('Search trips and expenses...')).toBeInTheDocument()
  })

  it('displays search value from filters', () => {
    const filtersWithSearch = { ...defaultFilters, search: 'test search' }
    
    render(
      <SearchFilter
        filters={filtersWithSearch}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByDisplayValue('test search')).toBeInTheDocument()
  })

  it('calls onFilterChange when search input changes', async () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    const searchInput = screen.getByPlaceholderText('Search trips and expenses...')
    fireEvent.change(searchInput, { target: { value: 'new search' } })

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        search: 'new search'
      })
    })
  })

  it('shows clear button when search has value', () => {
    const filtersWithSearch = { ...defaultFilters, search: 'test search' }
    
    render(
      <SearchFilter
        filters={filtersWithSearch}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByRole('button', { name: /clear/i })).toBeInTheDocument()
  })

  it('clears search when clear button is clicked', async () => {
    const filtersWithSearch = { ...defaultFilters, search: 'test search' }
    
    render(
      <SearchFilter
        filters={filtersWithSearch}
        onFilterChange={mockOnFilterChange}
      />
    )

    const clearButton = screen.getByRole('button', { name: /clear/i })
    fireEvent.click(clearButton)

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        search: ''
      })
    })
  })

  it('renders advanced filters button', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByText('Advanced Filters')).toBeInTheDocument()
  })

  it('toggles advanced filters when button is clicked', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    const advancedButton = screen.getByText('Advanced Filters')
    
    // Initially advanced filters should be hidden
    expect(screen.queryByText('Date From')).not.toBeInTheDocument()
    
    // Click to show advanced filters
    fireEvent.click(advancedButton)
    expect(screen.getByText('Date From')).toBeInTheDocument()
    
    // Click to hide advanced filters
    fireEvent.click(advancedButton)
    expect(screen.queryByText('Date From')).not.toBeInTheDocument()
  })

  it('shows advanced filters by default when showAdvanced is true', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        showAdvanced={true}
      />
    )

    expect(screen.getByText('Date From')).toBeInTheDocument()
    expect(screen.getByText('Date To')).toBeInTheDocument()
    expect(screen.getByText('Min Amount')).toBeInTheDocument()
    expect(screen.getByText('Max Amount')).toBeInTheDocument()
    expect(screen.getByText('Trip Type')).toBeInTheDocument()
    expect(screen.getByText('Category')).toBeInTheDocument()
  })

  it('calls onFilterChange when date filters change', async () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        showAdvanced={true}
      />
    )

    const dateFromInput = screen.getByLabelText('Date From')
    fireEvent.change(dateFromInput, { target: { value: '2024-01-01' } })

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        dateFrom: '2024-01-01'
      })
    })
  })

  it('calls onFilterChange when amount filters change', async () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        showAdvanced={true}
      />
    )

    const minAmountInput = screen.getByLabelText('Min Amount')
    fireEvent.change(minAmountInput, { target: { value: '10.50' } })

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        minAmount: '10.50'
      })
    })
  })

  it('calls onFilterChange when trip type changes', async () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        showAdvanced={true}
      />
    )

    const typeSelect = screen.getByLabelText('Trip Type')
    fireEvent.change(typeSelect, { target: { value: 'single' } })

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        type: 'single'
      })
    })
  })

  it('calls onFilterChange when category changes', async () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
        showAdvanced={true}
      />
    )

    const categoryInput = screen.getByLabelText('Category')
    fireEvent.change(categoryInput, { target: { value: 'food' } })

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...defaultFilters,
        category: 'food'
      })
    })
  })

  it('shows clear all button when filters are active', () => {
    const activeFilters = {
      ...defaultFilters,
      search: 'test',
      dateFrom: '2024-01-01',
      type: 'single'
    }
    
    render(
      <SearchFilter
        filters={activeFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByText('Clear All')).toBeInTheDocument()
  })

  it('clears all filters when clear all button is clicked', async () => {
    const activeFilters = {
      ...defaultFilters,
      search: 'test',
      dateFrom: '2024-01-01',
      type: 'single'
    }
    
    render(
      <SearchFilter
        filters={activeFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    const clearAllButton = screen.getByText('Clear All')
    fireEvent.click(clearAllButton)

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith(defaultFilters)
    })
  })

  it('displays active filter tags', () => {
    const activeFilters = {
      ...defaultFilters,
      search: 'test search',
      dateFrom: '2024-01-01',
      type: 'single'
    }
    
    render(
      <SearchFilter
        filters={activeFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.getByText('Search: test search')).toBeInTheDocument()
    expect(screen.getByText('From: 2024-01-01')).toBeInTheDocument()
    expect(screen.getByText('Type: single')).toBeInTheDocument()
  })

  it('removes individual filter when tag clear button is clicked', async () => {
    const activeFilters = {
      ...defaultFilters,
      search: 'test search',
      dateFrom: '2024-01-01'
    }
    
    render(
      <SearchFilter
        filters={activeFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    const searchTag = screen.getByText('Search: test search')
    const clearButton = searchTag.parentElement?.querySelector('button')
    
    if (clearButton) {
      fireEvent.click(clearButton)
    }

    await waitFor(() => {
      expect(mockOnFilterChange).toHaveBeenCalledWith({
        ...activeFilters,
        search: ''
      })
    })
  })

  it('does not show clear all button when no filters are active', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.queryByText('Clear All')).not.toBeInTheDocument()
  })

  it('does not show active filter tags when no filters are active', () => {
    render(
      <SearchFilter
        filters={defaultFilters}
        onFilterChange={mockOnFilterChange}
      />
    )

    expect(screen.queryByText(/Search:/)).not.toBeInTheDocument()
    expect(screen.queryByText(/From:/)).not.toBeInTheDocument()
    expect(screen.queryByText(/Type:/)).not.toBeInTheDocument()
  })
}) 