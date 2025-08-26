'use client'

import { useEffect } from 'react'
import { useAuthStore } from '@/store/auth'

interface AuthProviderProps {
  children: React.ReactNode
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const { initializeAuth } = useAuthStore()

  useEffect(() => {
    // Initialize auth state when the app loads
    initializeAuth()
  }, [initializeAuth])

  return <>{children}</>
}
