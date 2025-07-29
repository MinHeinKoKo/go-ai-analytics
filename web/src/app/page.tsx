'use client'

import { useRequireAuth } from '@/hooks/useAuth'
import Layout from '@/components/Layout'
import Dashboard from '@/components/Dashboard'

export default function Home() {
  const { isReady } = useRequireAuth()

  // Show loading spinner while checking auth state
  if (!isReady) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <Layout>
      <Dashboard />
    </Layout>
  )
}
