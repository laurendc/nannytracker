import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Trips from './pages/Trips'
import Expenses from './pages/Expenses'
import Summaries from './pages/Summaries'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/trips" element={<Trips />} />
            <Route path="/expenses" element={<Expenses />} />
            <Route path="/summaries" element={<Summaries />} />
          </Routes>
        </Layout>
      </Router>
    </QueryClientProvider>
  )
}

export default App 