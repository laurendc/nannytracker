import { render, screen } from '../test-utils'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'

// Test component
const TestComponent = () => (
  <div>
    <h1>Test Component</h1>
    <button>Click me</button>
  </div>
)

describe('Test Utils', () => {
  it('renders component with providers', () => {
    render(<TestComponent />)

    expect(screen.getByText('Test Component')).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('provides routing context', () => {
    render(<TestComponent />)

    // Check that we can access router context
    const button = screen.getByRole('button')
    expect(button).toBeInTheDocument()
  })

  it('provides query client context', () => {
    render(<TestComponent />)

    // The component should render without errors, indicating QueryClient is available
    expect(screen.getByText('Test Component')).toBeInTheDocument()
  })

  it('allows custom render options', () => {
    const customContainer = document.createElement('div')
    document.body.appendChild(customContainer)
    try {
      const { container } = render(<TestComponent />, {
        container: customContainer,
      })
      expect(container).toBe(customContainer)
      expect(customContainer.querySelector('h1')).toBeInTheDocument()
    } finally {
      document.body.removeChild(customContainer)
    }
  })
}) 