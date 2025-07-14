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

    const nannyTrackerElements = screen.getAllByText('NannyTracker')
    expect(nannyTrackerElements.length).toBeGreaterThan(0)
  })

  it('renders navigation links', () => {
    render(
      <TestWrapper>
        <Layout>
          <div>Test content</div>
        </Layout>
      </TestWrapper>
    )

    expect(screen.getAllByText('Dashboard').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Trips').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Expenses').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Summaries').length).toBeGreaterThan(0)
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

    const navElements = screen.getAllByRole('navigation')
    expect(navElements.length).toBeGreaterThan(0)

    const links = screen.getAllByRole('link')
    expect(links.length).toBeGreaterThan(0) // Has navigation links
  })
}) 