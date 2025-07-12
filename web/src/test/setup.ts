import '@testing-library/jest-dom'
import { vi } from 'vitest'
import { server } from './mocks/server'

// Mock ResizeObserver for Recharts
declare global {
  interface Window {
    ResizeObserver: typeof ResizeObserver
    confirm: typeof confirm
  }
}

window.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
})) as any

// Mock window.confirm for tests
window.confirm = vi.fn().mockReturnValue(true)

// Establish API mocking before all tests
beforeAll(() => server.listen())

// Reset any request handlers that we may add during the tests,
// so they don't affect other tests
afterEach(() => {
  server.resetHandlers()
  // Reset window.confirm mock after each test
  vi.mocked(window.confirm).mockClear()
  vi.mocked(window.confirm).mockReturnValue(true)
})

// Clean up after the tests are finished
afterAll(() => server.close()) 