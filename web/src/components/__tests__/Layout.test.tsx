import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Layout from '../Layout'

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
      <BrowserRouter>
        {children}
      </BrowserRouter>
    </QueryClientProvider>
  )
}

describe('Layout', () => {
  it('renders the header with app title', () => {
    render(
      <TestWrapper>
        <Layout>
          <div>Test content</div>
        </Layout>
      </TestWrapper>
    )

    expect(screen.getByText('NannyTracker')).toBeInTheDocument()
  })

  it('renders navigation links', () => {
    render(
      <TestWrapper>
        <Layout>
          <div>Test content</div>
        </Layout>
      </TestWrapper>
    )

    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Trips')).toBeInTheDocument()
    expect(screen.getByText('Expenses')).toBeInTheDocument()
    expect(screen.getByText('Summaries')).toBeInTheDocument()
  })

  it('renders children content', () => {
    render(
      <TestWrapper>
        <Layout>
          <div data-testid="test-content">Test content</div>
        </Layout>
      </TestWrapper>
    )

    expect(screen.getByTestId('test-content')).toBeInTheDocument()
  })

  it('has correct navigation structure', () => {
    render(
      <TestWrapper>
        <Layout>
          <div>Test content</div>
        </Layout>
      </TestWrapper>
    )

    const nav = screen.getByRole('navigation')
    expect(nav).toBeInTheDocument()

    const links = screen.getAllByRole('link')
    expect(links).toHaveLength(4) // Dashboard, Trips, Expenses, Summaries
  })
}) 