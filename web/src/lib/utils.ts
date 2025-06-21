/**
 * Utility functions for formatting and calculations
 */

/**
 * Format a number as currency
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(amount)
}

/**
 * Format a date string to a readable format
 */
export function formatDate(dateString: string): string {
  try {
    let date: Date
    
    // Handle different date formats
    if (dateString.includes('T') || dateString.includes(' ')) {
      // ISO string or datetime string
      date = new Date(dateString)
    } else {
      // Simple date string (YYYY-MM-DD)
      const [year, month, day] = dateString.split('-').map(Number)
      date = new Date(year, month - 1, day) // month is 0-indexed
    }
    
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  } catch {
    return 'Invalid Date'
  }
}

/**
 * Format miles with proper pluralization
 */
export function formatMiles(miles: number): string {
  const rounded = Math.round(miles * 10) / 10
  return rounded === 1 ? '1 mile' : `${rounded} miles`
}

/**
 * Calculate reimbursement amount based on miles
 * Current rate: $0.14 per mile
 */
export function calculateReimbursement(miles: number): number {
  const rate = 0.14
  return Math.round(miles * rate * 10) / 10
} 