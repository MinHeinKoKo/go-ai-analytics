import { create } from "zustand";
import { persist } from "zustand/middleware";

interface User {
  id: string;
  email: string;
  name: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  setUser: (user: User) => void;
  setLoading: (loading: boolean) => void;
  initializeAuth: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,
      login: (token: string, user: User) => {
        // Store token in localStorage for API requests
        if (typeof window !== "undefined") {
          localStorage.setItem("auth_token", token);
        }
        set({
          token,
          user,
          isAuthenticated: true,
          isLoading: false,
        });
      },
      logout: () => {
        // Clear token from localStorage
        if (typeof window !== "undefined") {
          localStorage.removeItem("auth_token");
        }
        set({
          token: null,
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },
      setUser: (user: User) => set({ user }),
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      initializeAuth: () => {
        // Initialize auth state from persisted storage
        const state = get();
        if (state.token && state.user) {
          // Ensure localStorage is in sync
          if (typeof window !== "undefined") {
            localStorage.setItem("auth_token", state.token);
          }
          set({ isAuthenticated: true, isLoading: false });
        } else {
          set({ isLoading: false });
        }
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Called when the persisted state is restored
        if (state) {
          state.initializeAuth();
        }
      },
    }
  )
);
