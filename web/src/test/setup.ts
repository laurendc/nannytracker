import '@testing-library/jest-dom'
import { server } from './mocks/server'

// Mock ResizeObserver for Recharts
declare global {
  interface Window {
    ResizeObserver: typeof ResizeObserver
  }
}

window.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
})) as any

// Establish API mocking before all tests
beforeAll(() => server.listen())

// Reset any request handlers that we may add during the tests,
// so they don't affect other tests
afterEach(() => server.resetHandlers())

// Clean up after the tests are finished
afterAll(() => server.close()) 