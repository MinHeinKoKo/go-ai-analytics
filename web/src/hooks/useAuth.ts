"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";

export function useAuth(redirectTo = "/login") {
  const { isAuthenticated, isLoading, user, token, initializeAuth } =
    useAuthStore();
  const router = useRouter();

  useEffect(() => {
    // Initialize auth state on mount
    initializeAuth();
  }, [initializeAuth]);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push(redirectTo);
    }
  }, [isAuthenticated, isLoading, router, redirectTo]);

  return {
    isAuthenticated,
    isLoading,
    user,
    token,
  };
}

export function useRequireAuth() {
  const auth = useAuth();

  if (auth.isLoading) {
    return { ...auth, isReady: false };
  }

  if (!auth.isAuthenticated) {
    return { ...auth, isReady: false };
  }

  return { ...auth, isReady: true };
}
