import { formatCurrency, formatDate, formatMiles, calculateReimbursement } from '../utils'

describe('Utils', () => {
  describe('formatCurrency', () => {
    it('formats currency correctly', () => {
      expect(formatCurrency(15.50)).toBe('$15.50')
      expect(formatCurrency(0)).toBe('$0.00')
      expect(formatCurrency(1234.56)).toBe('$1,234.56')
      expect(formatCurrency(0.99)).toBe('$0.99')
    })

    it('handles negative values', () => {
      expect(formatCurrency(-15.50)).toBe('-$15.50')
      expect(formatCurrency(-1234.56)).toBe('-$1,234.56')
    })

    it('handles zero', () => {
      expect(formatCurrency(0)).toBe('$0.00')
    })
  })

  describe('formatDate', () => {
    it('formats date correctly', () => {
      expect(formatDate('2024-12-18')).toBe('Dec 18, 2024')
      expect(formatDate('2024-01-01')).toBe('Jan 1, 2024')
      expect(formatDate('2024-12-31')).toBe('Dec 31, 2024')
    })

    it('handles different date formats', () => {
      expect(formatDate('2024-12-18T10:30:00Z')).toBe('Dec 18, 2024')
      expect(formatDate('2024-12-18 10:30:00')).toBe('Dec 18, 2024')
    })

    it('handles invalid dates gracefully', () => {
      expect(formatDate('invalid-date')).toBe('Invalid Date')
    })
  })

  describe('formatMiles', () => {
    it('formats miles correctly', () => {
      expect(formatMiles(5.2)).toBe('5.2 miles')
      expect(formatMiles(0)).toBe('0 miles')
      expect(formatMiles(1)).toBe('1 mile')
      expect(formatMiles(1.5)).toBe('1.5 miles')
    })

    it('handles decimal places correctly', () => {
      expect(formatMiles(5.25)).toBe('5.3 miles')
      expect(formatMiles(5.24)).toBe('5.2 miles')
      expect(formatMiles(5.26)).toBe('5.3 miles')
    })

    it('handles zero', () => {
      expect(formatMiles(0)).toBe('0 miles')
    })

    it('handles singular form', () => {
      expect(formatMiles(1)).toBe('1 mile')
      expect(formatMiles(1.0)).toBe('1 mile')
    })
  })

  describe('calculateReimbursement', () => {
    it('calculates reimbursement correctly', () => {
      expect(calculateReimbursement(5.2)).toBe(0.7) // 5.2 * 0.14 = 0.728, rounded to 0.7
      expect(calculateReimbursement(10)).toBe(1.4) // 10 * 0.14 = 1.4
      expect(calculateReimbursement(0)).toBe(0)
    })

    it('handles decimal miles correctly', () => {
      expect(calculateReimbursement(5.25)).toBe(0.7) // 5.25 * 0.14 = 0.735, rounded to 0.7
      expect(calculateReimbursement(5.26)).toBe(0.7) // 5.26 * 0.14 = 0.736, rounded to 0.7
    })

    it('handles zero miles', () => {
      expect(calculateReimbursement(0)).toBe(0)
    })

    it('handles large mileages', () => {
      expect(calculateReimbursement(100)).toBe(14) // 100 * 0.14 = 14
      expect(calculateReimbursement(1000)).toBe(140) // 1000 * 0.14 = 140
    })
  })
}) 