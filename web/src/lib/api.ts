import axios from 'axios'
import type { Trip, Expense, WeeklySummary, TripsResponse, ExpensesResponse, SummariesResponse } from '../types'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Trips API
export const tripsApi = {
  getAll: async (): Promise<Trip[]> => {
    const response = await api.get<TripsResponse>('/trips')
    return response.data.trips
  },
  
  create: async (trip: Omit<Trip, 'miles'>): Promise<Trip> => {
    const response = await api.post<Trip>('/trips', trip)
    return response.data
  },
  
  update: async (id: string, trip: Trip): Promise<Trip> => {
    const response = await api.put<Trip>(`/trips/${id}`, trip)
    return response.data
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/trips/${id}`)
  },
}

// Expenses API
export const expensesApi = {
  getAll: async (): Promise<Expense[]> => {
    const response = await api.get<ExpensesResponse>('/expenses')
    return response.data.expenses
  },
  
  create: async (expense: Expense): Promise<Expense> => {
    const response = await api.post<Expense>('/expenses', expense)
    return response.data
  },
  
  update: async (id: string, expense: Expense): Promise<Expense> => {
    const response = await api.put<Expense>(`/expenses/${id}`, expense)
    return response.data
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/expenses/${id}`)
  },
}

// Summaries API
export const summariesApi = {
  getAll: async (): Promise<WeeklySummary[]> => {
    const response = await api.get<SummariesResponse>('/summaries')
    return response.data.summaries
  },
}

// Health check
export const healthApi = {
  check: async (): Promise<{ status: string; service: string }> => {
    const response = await api.get('/health')
    return response.data
  },
} 