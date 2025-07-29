'use client'

import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { analyticsApi } from '@/lib/api'
import { formatCurrency } from '@/lib/utils'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line, PieChart, Pie, Cell } from 'recharts'
import { Users, ShoppingCart, DollarSign, TrendingUp, Target, Brain } from 'lucide-react'

export default function Dashboard() {
  const { data: dashboardData, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: () => analyticsApi.getDashboard(),
  })

  const dashboard = dashboardData?.data.dashboard

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-gray-900"></div>
      </div>
    )
  }

  const metrics = [
    {
      title: 'Total Customers',
      value: dashboard?.total_customers || 0,
      icon: Users,
      description: 'Active customer base',
    },
    {
      title: 'Total Revenue',
      value: formatCurrency(dashboard?.total_revenue || 0),
      icon: DollarSign,
      description: 'Total revenue generated',
    },
    {
      title: 'Total Purchases',
      value: dashboard?.total_purchases || 0,
      icon: ShoppingCart,
      description: 'Total transactions',
    },
    {
      title: 'Avg Order Value',
      value: formatCurrency(dashboard?.avg_order_value || 0),
      icon: TrendingUp,
      description: 'Average order value',
    },
    {
      title: 'Active Campaigns',
      value: dashboard?.active_campaigns || 0,
      icon: Target,
      description: 'Currently running campaigns',
    },
    {
      title: 'Total Campaigns',
      value: dashboard?.total_campaigns || 0,
      icon: Brain,
      description: 'All marketing campaigns',
    },
  ]

  // Sample data for charts
  const revenueData = [
    { month: 'Jan', revenue: 12000, customers: 120 },
    { month: 'Feb', revenue: 15000, customers: 150 },
    { month: 'Mar', revenue: 18000, customers: 180 },
    { month: 'Apr', revenue: 22000, customers: 220 },
    { month: 'May', revenue: 25000, customers: 250 },
    { month: 'Jun', revenue: 28000, customers: 280 },
  ]

  const segmentData = [
    { name: 'High Value', value: 30, color: '#8884d8' },
    { name: 'Medium Value', value: 45, color: '#82ca9d' },
    { name: 'Low Value', value: 25, color: '#ffc658' },
  ]

  const campaignData = [
    { name: 'Email', performance: 85, cost: 2000 },
    { name: 'Social Media', performance: 92, cost: 3500 },
    { name: 'Display Ads', performance: 78, cost: 4000 },
    { name: 'Search Ads', performance: 88, cost: 5000 },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-primary">AI Analytics Dashboard</h1>
        <p className="text-muted-foreground">
          Marketing & targeting insights powered by artificial intelligence
        </p>
      </div>

      {/* Metrics Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {metrics.map((metric, index) => {
          const Icon = metric.icon
          return (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {metric.title}
                </CardTitle>
                <Icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metric.value}</div>
                <p className="text-xs text-muted-foreground">
                  {metric.description}
                </p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Charts */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="segments">Customer Segments</TabsTrigger>
          <TabsTrigger value="campaigns">Campaign Performance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Revenue Trend</CardTitle>
                <CardDescription>
                  Monthly revenue and customer growth
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={revenueData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="month" />
                    <YAxis />
                    <Tooltip />
                    <Line
                      type="monotone"
                      dataKey="revenue"
                      stroke="#8884d8"
                      strokeWidth={2}
                      name="Revenue ($)"
                    />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Customer Growth</CardTitle>
                <CardDescription>
                  New customers acquired monthly
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={revenueData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="month" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="customers" fill="#82ca9d" name="New Customers" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="segments" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Customer Segments</CardTitle>
                <CardDescription>
                  AI-generated customer segmentation
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={segmentData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {segmentData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Segment Insights</CardTitle>
                <CardDescription>
                  Key characteristics of each segment
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-[#8884d8] rounded-full mr-2"></div>
                    <span className="text-sm font-medium">High Value (30%)</span>
                  </div>
                  <p className="text-xs text-muted-foreground ml-5">
                    High spending, frequent purchases, premium products
                  </p>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-[#82ca9d] rounded-full mr-2"></div>
                    <span className="text-sm font-medium">Medium Value (45%)</span>
                  </div>
                  <p className="text-xs text-muted-foreground ml-5">
                    Moderate spending, regular purchases, price-conscious
                  </p>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-[#ffc658] rounded-full mr-2"></div>
                    <span className="text-sm font-medium">Low Value (25%)</span>
                  </div>
                  <p className="text-xs text-muted-foreground ml-5">
                    Low spending, infrequent purchases, deal-seekers
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="campaigns" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Campaign Performance</CardTitle>
              <CardDescription>
                AI-optimized marketing campaign results
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={campaignData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="performance" fill="#8884d8" name="Performance Score" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
