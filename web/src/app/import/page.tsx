'use client'

import { useRequireAuth } from '@/hooks/useAuth'
import Layout from '@/components/Layout'
import DataImport from '@/components/DataImport'

export default function ImportPage() {
  const { isReady } = useRequireAuth()

  if (!isReady) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <Layout>
      <DataImport />
    </Layout>
  )
}
