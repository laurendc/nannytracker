import React from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line } from 'recharts'

interface ChartData {
  week: string
  miles: number
  amount: number
  expenses: number
}

interface ChartsProps {
  chartData: ChartData[]
}

const Charts: React.FC<ChartsProps> = ({ chartData }) => {
  return (
    <div className="space-y-6">
      {/* Weekly Miles Chart */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Weekly Miles</h2>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="week" />
            <YAxis />
            <Tooltip formatter={(value) => [`${value} miles`, 'Miles']} />
            <Bar dataKey="miles" fill="#3b82f6" />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Weekly Amount Chart */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Weekly Reimbursement Amount</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="week" />
            <YAxis />
            <Tooltip formatter={(value) => [`$${value}`, 'Amount']} />
            <Line type="monotone" dataKey="amount" stroke="#10b981" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

export default Charts 