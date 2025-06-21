import { render, screen } from '@testing-library/react'
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

    expect(screen.getByText('NannyTracker')).toBeInTheDocument()
  })

  it('renders navigation links', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    expect(screen.getAllByText('Dashboard').length).toBeGreaterThan(0)
    expect(screen.getByText('Trips')).toBeInTheDocument()
    expect(screen.getByText('Expenses')).toBeInTheDocument()
    expect(screen.getByText('Summaries')).toBeInTheDocument()
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

    // Check for header
    expect(screen.getByRole('banner')).toBeInTheDocument()
    
    // Check for navigation
    expect(screen.getByRole('navigation')).toBeInTheDocument()
    
    // Check for main content
    expect(screen.getByRole('main')).toBeInTheDocument()
  })

  it('renders dashboard by default', () => {
    render(
      <TestWrapper>
        <App />
      </TestWrapper>
    )

    // Dashboard should be the default route
    expect(screen.getByText(/Overview of your mileage and expense tracking/)).toBeInTheDocument()
  })
}) 