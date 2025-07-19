import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from '../../App'

// Create a simple wrapper for testing
const TestWrapper = ({ children }: { children: React.ReactNode }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}

describe('App', () => {
  it('renders the app title', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    expect(screen.getAllByText('NannyTracker').length).toBeGreaterThan(0)
  })

  it('renders navigation links', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    expect(screen.getAllByText('Dashboard').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Trips').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Expenses').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Summaries').length).toBeGreaterThan(0)
  })

  it('renders main content area', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    // Check that main content is rendered
    expect(screen.getByRole('main')).toBeInTheDocument()
  })

  it('has proper app structure', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    // Check for headers (mobile and desktop)
    expect(screen.getAllByRole('banner').length).toBeGreaterThan(0)
    
    // Check for navigation (mobile and desktop)
    expect(screen.getAllByRole('navigation').length).toBeGreaterThan(0)
    
    // Check for main content
    expect(screen.getByRole('main')).toBeInTheDocument()
  })

  it('renders dashboard by default', async () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    // Dashboard should be the default route - wait for it to load
    await waitFor(() => {
      expect(screen.getByText(/Overview of your mileage and expense tracking/)).toBeInTheDocument()
    })
  })
}) 