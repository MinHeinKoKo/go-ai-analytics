'use client'

import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Upload, Download, FileText, AlertCircle, CheckCircle, Info } from 'lucide-react'

// API functions (inline to avoid import issues)
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface ImportResult {
  success_count: number
  total_rows: number
  imported: number
  errors?: string[]
}

// API functions
const getImportTemplates = async () => {
  const response = await fetch(`${API_BASE_URL}/api/v1/import/templates`)
  return response.json()
}

const importFile = async (type: string, file: File) => {
  const formData = new FormData()
  formData.append('file', file)

  const token = localStorage.getItem('auth_token')
  const response = await fetch(`${API_BASE_URL}/api/v1/import/${type}`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  })

  return response.json()
}

export default function DataImport() {
  const [selectedFiles, setSelectedFiles] = useState<{ [key: string]: File | null }>({
    customers: null,
    purchases: null,
    campaigns: null,
    performance: null,
  })
  const [importResults, setImportResults] = useState<{ [key: string]: ImportResult }>({})

  // Get import templates
  const { data: templatesData, isLoading } = useQuery({
    queryKey: ['import-templates'],
    queryFn: getImportTemplates,
  })

  // Import mutation
  const importMutation = useMutation({
    mutationFn: ({ type, file }: { type: string, file: File }) => importFile(type, file),
    onSuccess: (data, variables) => {
      setImportResults(prev => ({ ...prev, [variables.type]: data }))
    },
  })

  const handleFileSelect = (type: string, file: File | null) => {
    setSelectedFiles(prev => ({ ...prev, [type]: file }))
    // Clear previous results when new file is selected
    if (importResults[type]) {
      setImportResults(prev => {
        const newResults = { ...prev }
        delete newResults[type]
        return newResults
      })
    }
  }

  const handleImport = (type: string) => {
    const file = selectedFiles[type]
    if (!file) return
    importMutation.mutate({ type, file })
  }

  const downloadSample = (type: string) => {
    const link = document.createElement('a')
    link.href = `${API_BASE_URL}/api/v1/import/sample/${type}`
    link.download = `sample_${type}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-gray-900"></div>
      </div>
    )
  }

  const templates = templatesData?.templates || {
    customers: {
      required_headers: ['customer_id', 'age', 'gender', 'location', 'income_range', 'registration_date', 'preferred_category'],
      data_types: {
        customer_id: 'string (unique identifier)',
        age: 'integer (18-100)',
        gender: 'string (Male/Female/Other)',
        location: 'string (city/state)',
        income_range: 'string (e.g., $50k-$75k)',
        registration_date: 'date (YYYY-MM-DD format)',
        preferred_category: 'string (product category)'
      },
      example_row: 'CUST00001,25,Female,New York,$50k-$75k,2024-01-15,Fashion'
    },
    purchases: {
      required_headers: ['customer_id', 'product_id', 'category', 'amount', 'quantity', 'purchase_date', 'channel'],
      data_types: {
        customer_id: 'string (must exist in customers)',
        product_id: 'string (product identifier)',
        category: 'string (product category)',
        amount: 'decimal (purchase amount)',
        quantity: 'integer (number of items)',
        purchase_date: 'date (YYYY-MM-DD format)',
        channel: 'string (online/store)'
      },
      example_row: 'CUST00001,PROD001,Fashion,89.99,1,2024-01-20,online'
    },
    campaigns: {
      required_headers: ['campaign_id', 'name', 'type', 'target_segment', 'budget', 'start_date', 'end_date', 'status'],
      data_types: {
        campaign_id: 'string (unique identifier)',
        name: 'string (campaign name)',
        type: 'string (email/social/display/search)',
        target_segment: 'string (target audience)',
        budget: 'decimal (campaign budget)',
        start_date: 'date (YYYY-MM-DD format)',
        end_date: 'date (YYYY-MM-DD format)',
        status: 'string (active/paused/completed)'
      },
      example_row: 'CAMP0001,Summer Sale,email,Fashion Lovers,5000.00,2024-06-01,2024-06-30,completed'
    },
    performance: {
      required_headers: ['campaign_id', 'impressions', 'clicks', 'conversions', 'revenue', 'cost', 'date'],
      data_types: {
        campaign_id: 'string (must exist in campaigns)',
        impressions: 'integer (ad impressions)',
        clicks: 'integer (ad clicks)',
        conversions: 'integer (conversions)',
        revenue: 'decimal (revenue generated)',
        cost: 'decimal (campaign cost)',
        date: 'date (YYYY-MM-DD format)'
      },
      example_row: 'CAMP0001,10000,500,25,2500.00,1000.00,2024-06-01'
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-primary">Data Import</h1>
        <p className="text-muted-foreground">
          Import your data using CSV files. Download sample files to see the required format.
        </p>
      </div>

      {/* General Guidelines */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Info className="h-5 w-5" />
            Import Guidelines
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <h4 className="font-medium mb-2">File Requirements:</h4>
              <ul className="text-sm text-muted-foreground space-y-1">
                <li>• CSV format only</li>
                <li>• Headers must be in the first row</li>
                <li>• Maximum file size: 10MB</li>
                <li>• Maximum 10,000 rows per import</li>
              </ul>
            </div>
            <div>
              <h4 className="font-medium mb-2">Data Format:</h4>
              <ul className="text-sm text-muted-foreground space-y-1">
                <li>• Dates: YYYY-MM-DD format</li>
                <li>• Decimals: Use dot (.) separator</li>
                <li>• Text: No commas in text fields</li>
                <li>• IDs: Must be unique where specified</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="customers" className="space-y-4">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="customers">Customers</TabsTrigger>
          <TabsTrigger value="purchases">Purchases</TabsTrigger>
          <TabsTrigger value="campaigns">Campaigns</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>

        {/* Customers Tab */}
        <TabsContent value="customers" className="space-y-4">
          <div className="grid gap-6 md:grid-cols-2">
            {/* Template Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Customers Template
                </CardTitle>
                <CardDescription>
                  Required format and data types for customers import
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h4 className="font-medium mb-2">Required Headers:</h4>
                  <div className="flex flex-wrap gap-1">
                    {templates.customers.required_headers.map((header: string) => (
                      <span
                        key={header}
                        className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded"
                      >
                        {header}
                      </span>
                    ))}
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Data Types:</h4>
                  <div className="space-y-1 text-sm">
                    {Object.entries(templates.customers.data_types).map(([field, type]: [string, any]) => (
                      <div key={field} className="flex justify-between">
                        <span className="font-mono text-xs">{field}:</span>
                        <span className="text-muted-foreground text-xs">{type}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Example Row:</h4>
                  <code className="text-xs bg-gray-100 p-2 rounded block break-all">
                    {templates.customers.example_row}
                  </code>
                </div>

                <Button
                  onClick={() => downloadSample('customers')}
                  variant="outline"
                  className="w-full"
                >
                  <Download className="h-4 w-4 mr-2" />
                  Download Sample CSV
                </Button>
              </CardContent>
            </Card>

            {/* Import Section */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Upload className="h-5 w-5" />
                  Import Customers
                </CardTitle>
                <CardDescription>
                  Upload your CSV file to import customers data
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <input
                    type="file"
                    accept=".csv"
                    onChange={(e) => handleFileSelect('customers', e.target.files?.[0] || null)}
                    className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                  />
                </div>

                {selectedFiles.customers && (
                  <div className="p-3 bg-blue-50 rounded-lg">
                    <p className="text-sm font-medium">Selected File:</p>
                    <p className="text-sm text-muted-foreground">{selectedFiles.customers?.name}</p>
                    <p className="text-xs text-muted-foreground">
                      Size: {((selectedFiles.customers?.size || 0) / 1024).toFixed(1)} KB
                    </p>
                  </div>
                )}

                <Button
                  onClick={() => handleImport('customers')}
                  disabled={!selectedFiles.customers || importMutation.isPending}
                  className="w-full"
                >
                  {importMutation.isPending ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      Importing...
                    </>
                  ) : (
                    <>
                      <Upload className="h-4 w-4 mr-2" />
                      Import Data
                    </>
                  )}
                </Button>

                {/* Import Results */}
                {importResults.customers && (
                  <div className="space-y-3">
                    <div className="flex items-center gap-2 p-3 bg-green-50 rounded-lg">
                      <CheckCircle className="h-5 w-5 text-green-600" />
                      <div>
                        <p className="font-medium text-green-800">Import Completed</p>
                        <p className="text-sm text-green-600">
                          {importResults.customers.success_count} of {importResults.customers.total_rows} rows imported successfully
                        </p>
                      </div>
                    </div>

                    {importResults.customers.errors && importResults.customers.errors.length > 0 && (
                      <div className="p-3 bg-red-50 rounded-lg">
                        <div className="flex items-center gap-2 mb-2">
                          <AlertCircle className="h-5 w-5 text-red-600" />
                          <p className="font-medium text-red-800">Import Errors</p>
                        </div>
                        <div className="max-h-32 overflow-y-auto">
                          {importResults.customers.errors.map((error, index) => (
                            <p key={index} className="text-sm text-red-600 mb-1">
                              {error}
                            </p>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                )}

                {importMutation.isError && (
                  <div className="flex items-center gap-2 p-3 bg-red-50 rounded-lg">
                    <AlertCircle className="h-5 w-5 text-red-600" />
                    <div>
                      <p className="font-medium text-red-800">Import Failed</p>
                      <p className="text-sm text-red-600">
                        {(importMutation.error as any)?.message || 'An error occurred during import'}
                      </p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Purchases Tab */}
        <TabsContent value="purchases" className="space-y-4">
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Purchases Template
                </CardTitle>
                <CardDescription>
                  Required format and data types for purchases import
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h4 className="font-medium mb-2">Required Headers:</h4>
                  <div className="flex flex-wrap gap-1">
                    {templates.purchases.required_headers.map((header: string) => (
                      <span key={header} className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                        {header}
                      </span>
                    ))}
                  </div>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Example Row:</h4>
                  <code className="text-xs bg-gray-100 p-2 rounded block break-all">
                    {templates.purchases.example_row}
                  </code>
                </div>
                <Button onClick={() => downloadSample('purchases')} variant="outline" className="w-full">
                  <Download className="h-4 w-4 mr-2" />
                  Download Sample CSV
                </Button>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Upload className="h-5 w-5" />
                  Import Purchases
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => handleFileSelect('purchases', e.target.files?.[0] || null)}
                  className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                />
                {selectedFiles.purchases && (
                  <div className="p-3 bg-blue-50 rounded-lg">
                    <p className="text-sm font-medium">Selected: {selectedFiles.purchases?.name}</p>
                  </div>
                )}
                <Button
                  onClick={() => handleImport('purchases')}
                  disabled={!selectedFiles.purchases || importMutation.isPending}
                  className="w-full"
                >
                  <Upload className="h-4 w-4 mr-2" />
                  Import Data
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Campaigns Tab */}
        <TabsContent value="campaigns" className="space-y-4">
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Campaigns Template
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h4 className="font-medium mb-2">Required Headers:</h4>
                  <div className="flex flex-wrap gap-1">
                    {templates.campaigns.required_headers.map((header: string) => (
                      <span key={header} className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                        {header}
                      </span>
                    ))}
                  </div>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Example Row:</h4>
                  <code className="text-xs bg-gray-100 p-2 rounded block break-all">
                    {templates.campaigns.example_row}
                  </code>
                </div>
                <Button onClick={() => downloadSample('campaigns')} variant="outline" className="w-full">
                  <Download className="h-4 w-4 mr-2" />
                  Download Sample CSV
                </Button>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Upload className="h-5 w-5" />
                  Import Campaigns
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => handleFileSelect('campaigns', e.target.files?.[0] || null)}
                  className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                />
                {selectedFiles.campaigns && (
                  <div className="p-3 bg-blue-50 rounded-lg">
                    <p className="text-sm font-medium">Selected: {selectedFiles.campaigns?.name}</p>
                  </div>
                )}
                <Button
                  onClick={() => handleImport('campaigns')}
                  disabled={!selectedFiles.campaigns || importMutation.isPending}
                  className="w-full"
                >
                  <Upload className="h-4 w-4 mr-2" />
                  Import Data
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Performance Tab */}
        <TabsContent value="performance" className="space-y-4">
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Performance Template
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h4 className="font-medium mb-2">Required Headers:</h4>
                  <div className="flex flex-wrap gap-1">
                    {templates.performance.required_headers.map((header: string) => (
                      <span key={header} className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                        {header}
                      </span>
                    ))}
                  </div>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Example Row:</h4>
                  <code className="text-xs bg-gray-100 p-2 rounded block break-all">
                    {templates.performance.example_row}
                  </code>
                </div>
                <Button onClick={() => downloadSample('performance')} variant="outline" className="w-full">
                  <Download className="h-4 w-4 mr-2" />
                  Download Sample CSV
                </Button>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Upload className="h-5 w-5" />
                  Import Performance
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => handleFileSelect('performance', e.target.files?.[0] || null)}
                  className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                />
                {selectedFiles.performance && (
                  <div className="p-3 bg-blue-50 rounded-lg">
                    <p className="text-sm font-medium">Selected: {selectedFiles.performance?.name}</p>
                  </div>
                )}
                <Button
                  onClick={() => handleImport('performance')}
                  disabled={!selectedFiles.performance || importMutation.isPending}
                  className="w-full"
                >
                  <Upload className="h-4 w-4 mr-2" />
                  Import Data
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
