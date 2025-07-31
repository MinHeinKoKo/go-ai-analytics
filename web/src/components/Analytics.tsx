'use client'

import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { analyticsApi } from '@/lib/api'
import { Brain, Users, Target, TrendingUp, Zap } from 'lucide-react'

export default function Analytics() {
  const [selectedCustomer, setSelectedCustomer] = useState('')
  const [selectedCampaign, setSelectedCampaign] = useState('')
  const queryClient = useQueryClient()

  const { data: customersData } = useQuery({
    queryKey: ['customers'],
    queryFn: () => analyticsApi.getCustomers(100, 0),
  })

  const { data: campaignsData } = useQuery({
    queryKey: ['campaigns'],
    queryFn: () => analyticsApi.getCampaigns(),
  })

  const segmentationMutation = useMutation({
    mutationFn: analyticsApi.performSegmentation,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['segments'] })
    },
  })

  const predictionMutation = useMutation({
    mutationFn: analyticsApi.predictCustomerBehavior,
  })

  const optimizationMutation = useMutation({
    mutationFn: analyticsApi.optimizeCampaign,
  })

  const sampleDataMutation = useMutation({
    mutationFn: analyticsApi.generateSampleData,
    onSuccess: () => {
      queryClient.invalidateQueries()
    },
  })

  const handleSegmentation = () => {
    segmentationMutation.mutate({
      algorithm: 'kmeans',
      features: ['total_spent', 'purchase_frequency', 'age'],
      parameters: { clusters: 3 }
    })
  }

  const handlePrediction = () => {
    if (!selectedCustomer) return
    predictionMutation.mutate({
      customer_id: selectedCustomer,
      prediction_type: 'churn'
    })
  }

  const handleOptimization = () => {
    if (!selectedCampaign) return
    optimizationMutation.mutate({
      campaign_id: selectedCampaign,
      objective: 'maximize_roas'
    })
  }

  const customers = customersData?.data.customers || []
  const campaigns = campaignsData?.data.campaigns || []

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-primary">AI Analytics</h1>
        <p className="text-muted-foreground">
          Advanced AI-powered analytics for marketing optimization
        </p>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Zap className="h-5 w-5" />
            Quick Actions
          </CardTitle>
          <CardDescription>
            Generate sample data and run AI analytics
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button
            onClick={() => sampleDataMutation.mutate()}
            disabled={sampleDataMutation.isPending}
            className="mr-4"
          >
            {sampleDataMutation.isPending ? 'Generating...' : 'Generate Sample Data'}
          </Button>
          {sampleDataMutation.isSuccess && (
            <p className="text-sm text-green00 mt-2">
              Sample data generated successfully!
            </p>
          )}
        </CardContent>
      </Card>

      <Tabs defaultValue="segmentation" className="space-y-4">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="segmentation">Customer Segmentation</TabsTrigger>
          <TabsTrigger value="prediction">Behavior Prediction</TabsTrigger>
          <TabsTrigger value="optimization">Campaign Optimization</TabsTrigger>
        </TabsList>

        <TabsContent value="segmentation" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Customer Segmentation
              </CardTitle>
              <CardDescription>
                Use AI to automatically segment customers based on behavior patterns
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <h4 className="text-sm font-medium mb-2">Algorithm: K-Means Clustering</h4>
                  <p className="text-sm text-muted-foreground mb-4">
                    Groups customers into segments based on spending patterns, purchase frequency, and demographics.
                  </p>
                  <Button
                    onClick={handleSegmentation}
                    disabled={segmentationMutation.isPending}
                  >
                    {segmentationMutation?.isPending ? 'Analyzing...' : 'Run Segmentation'}
                  </Button>
                </div>
                <div>
                  <h4 className="text-sm font-medium mb-2">Features Used:</h4>
                  <ul className="text-sm text-muted-foreground space-y-1">
                    <li>• Total amount spent</li>
                    <li>• Purchase frequency</li>
                    <li>• Customer age</li>
                    <li>• Registration date</li>
                  </ul>
                </div>
              </div>

              {segmentationMutation.isSuccess && (
                <div className="mt-4 p-4 bg-green-50 rounded-lg">
                  <h4 className="text-sm font-medium text-green-800 mb-2">Segmentation Results:</h4>
                  <div className="space-y-2">
                    {segmentationMutation?.data?.data?.segments?.map((segment, index) => (
                      <div key={index} className="text-sm">
                        <span className="font-medium">{segment.name}</span>: {segment.size} customers
                        <p className="text-xs text-muted-foreground">{segment.description}</p>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="prediction" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Brain className="h-5 w-5" />
                Customer Behavior Prediction
              </CardTitle>
              <CardDescription>
                Predict customer churn, lifetime value, and next purchase timing
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="text-sm font-medium">Select Customer:</label>
                <select
                  className="w-full mt-1 p-2 border rounded-md"
                  value={selectedCustomer}
                  onChange={(e) => setSelectedCustomer(e.target.value)}
                >
                  <option value="">Choose a customer...</option>
                  {customers.map((customer) => (
                    <option key={customer.id} value={customer.customer_id}>
                      {customer.customer_id} - {customer.preferred_category}
                    </option>
                  ))}
                </select>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Churn Prediction</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    Likelihood of customer stopping purchases
                  </p>
                  <Button
                    size="sm"
                    onClick={handlePrediction}
                    disabled={!selectedCustomer || predictionMutation.isPending}
                  >
                    Predict Churn
                  </Button>
                </div>
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Lifetime Value</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    Predicted total customer value
                  </p>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled
                  >
                    Coming Soon
                  </Button>
                </div>
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Next Purchase</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    When customer will likely buy again
                  </p>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled
                  >
                    Coming Soon
                  </Button>
                </div>
              </div>

              {predictionMutation.isSuccess && (
                <div className="mt-4 p-4 bg-blue-50 rounded-lg">
                  <h4 className="text-sm font-medium text-blue-800 mb-2">Prediction Results:</h4>
                  <div className="text-sm">
                    <p><span className="font-medium">Customer:</span> {predictionMutation.data.data.prediction.customer_id}</p>
                    <p><span className="font-medium">Churn Probability:</span> {(predictionMutation.data.data.prediction.probability * 100).toFixed(1)}%</p>
                    <p><span className="font-medium">Confidence:</span> {(predictionMutation.data.data.prediction.confidence * 100).toFixed(1)}%</p>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="optimization" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Target className="h-5 w-5" />
                Campaign Optimization
              </CardTitle>
              <CardDescription>
                AI-powered recommendations to improve campaign performance
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="text-sm font-medium">Select Campaign:</label>
                <select
                  className="w-full mt-1 p-2 border rounded-md"
                  value={selectedCampaign}
                  onChange={(e) => setSelectedCampaign(e.target.value)}
                >
                  <option value="">Choose a campaign...</option>
                  {campaigns.map((campaign) => (
                    <option key={campaign.id} value={campaign.campaign_id}>
                      {campaign.name} - {campaign.type}
                    </option>
                  ))}
                </select>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Maximize ROAS</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    Optimize for return on ad spend
                  </p>
                  <Button
                    size="sm"
                    onClick={handleOptimization}
                    disabled={!selectedCampaign || optimizationMutation.isPending}
                  >
                    Optimize
                  </Button>
                </div>
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Minimize Cost</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    Reduce campaign spending
                  </p>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled
                  >
                    Coming Soon
                  </Button>
                </div>
                <div className="p-4 border rounded-lg">
                  <h4 className="text-sm font-medium mb-2">Max Conversions</h4>
                  <p className="text-xs text-muted-foreground mb-3">
                    Increase conversion rate
                  </p>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled
                  >
                    Coming Soon
                  </Button>
                </div>
              </div>

              {optimizationMutation.isSuccess && (
                <div className="mt-4 p-4 bg-purple-50 rounded-lg">
                  <h4 className="text-sm font-medium text-purple-800 mb-2">Optimization Recommendations:</h4>
                  <div className="text-sm space-y-2">
                    <div>
                      <span className="font-medium">Optimization Score:</span> {optimizationMutation.data.data.optimization.optimization_score}/100
                    </div>
                    <div>
                      <span className="font-medium">Recommendations:</span>
                      <ul className="list-disc list-inside mt-1 text-xs">
                        {optimizationMutation.data.data.optimization.recommendations?.map((rec: string, index: number) => (
                          <li key={index}>{rec}</li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
