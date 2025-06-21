export interface Trip {
  date: string
  origin: string
  destination: string
  miles: number
  type: 'single' | 'round'
}

export interface Expense {
  date: string
  amount: number
  description: string
}

export interface WeeklySummary {
  weekStart: string
  weekEnd: string
  totalMiles: number
  totalAmount: number
  expenses: Expense[]
  trips: Trip[]
}

export interface TripTemplate {
  name: string
  origin: string
  destination: string
  tripType: 'single' | 'round'
  notes?: string
}

export interface RecurringTrip {
  startDate: string
  endDate: string
  weekday: number
  origin: string
  destination: string
  miles: number
  type: 'single' | 'round'
}

export interface ApiResponse<T> {
  data: T
  count: number
}

export interface TripsResponse extends ApiResponse<Trip[]> {
  trips: Trip[]
}

export interface ExpensesResponse extends ApiResponse<Expense[]> {
  expenses: Expense[]
}

export interface SummariesResponse extends ApiResponse<WeeklySummary[]> {
  summaries: WeeklySummary[]
} 