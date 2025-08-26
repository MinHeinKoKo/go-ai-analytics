# Authentication Implementation Guide

## Overview
The authentication system has been enhanced to properly store and manage user tokens and data across the application.

## Key Improvements

### 1. Enhanced Auth Store (`web/src/store/auth.ts`)
- **Persistent Storage**: Uses Zustand persist middleware to store auth state
- **Token Management**: Automatically syncs tokens with localStorage
- **Loading States**: Proper loading state management during auth checks
- **Initialization**: Proper auth state initialization on app load
- **SSR Safety**: Handles server-side rendering safely

### 2. Improved API Client (`web/src/lib/api.ts`)
- **Automatic Token Injection**: Adds Bearer token to all authenticated requests
- **Error Handling**: Automatically handles 401 errors and redirects to login
- **SSR Safety**: Only accesses localStorage on client side
- **Response Interceptors**: Handles token expiration gracefully

### 3. Enhanced Login Flow (`web/src/app/login/page.tsx`)
- **Better Error Handling**: Shows user-friendly error messages
- **Loading States**: Proper loading indicators during auth operations
- **Redirect Prevention**: Prevents rendering login form if already authenticated
- **Form Validation**: Client-side validation for better UX

### 4. Auth Hook (`web/src/hooks/useAuth.ts`)
- **Reusable Logic**: Centralized authentication logic
- **Route Protection**: Easy-to-use hook for protecting routes
- **Loading Management**: Handles loading states consistently

### 5. Auth Provider (`web/src/components/AuthProvider.tsx`)
- **Global Initialization**: Ensures auth state is initialized app-wide
- **Context Management**: Provides auth context to all components

## How Authentication Works

### 1. User Registration/Login
```typescript
// User submits login form
const loginMutation = useMutation({
  mutationFn: ({ email, password }) => authApi.login(email, password),
  onSuccess: (data) => {
    // Store token and user data
    login(data.data.token, {
      id: data.data.user.id,
      email: data.data.user.email,
      name: data.data.user.name,
    })
    router.push('/')
  }
})
```

### 2. Token Storage
```typescript
login: (token: string, user: User) => {
  // Store in localStorage for API requests
  if (typeof window !== 'undefined') {
    localStorage.setItem("auth_token", token);
  }
  // Update Zustand store (persisted automatically)
  set({ 
    token, 
    user, 
    isAuthenticated: true, 
    isLoading: false 
  });
}
```

### 3. API Request Authentication
```typescript
// Automatic token injection
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem("auth_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});
```

### 4. Route Protection
```typescript
// Using the auth hook
export default function ProtectedPage() {
  const { isReady } = useRequireAuth()

  if (!isReady) {
    return <LoadingSpinner />
  }

  return <ProtectedContent />
}
```

## Data Flow

1. **App Initialization**
   - AuthProvider initializes auth state from persisted storage
   - Token is restored to localStorage if available
   - User data is restored from Zustand persist

2. **Login Process**
   - User submits credentials
   - API returns token and user data
   - Auth store saves both token and user
   - Token is stored in localStorage for API requests
   - User is redirected to dashboard

3. **API Requests**
   - Request interceptor adds Bearer token automatically
   - If 401 error occurs, user is redirected to login
   - Token is cleared from storage on logout

4. **Route Protection**
   - Protected routes check authentication state
   - Loading spinner shown while checking auth
   - Redirect to login if not authenticated

## Storage Strategy

### localStorage
- Stores JWT token for API requests
- Cleared on logout or 401 errors
- Only accessed on client side (SSR safe)

### Zustand Persist
- Stores user data and authentication state
- Survives page refreshes and browser restarts
- Automatically syncs with localStorage

## Security Features

1. **Token Expiration Handling**: Automatic logout on 401 errors
2. **SSR Safety**: No localStorage access during server rendering
3. **Route Protection**: Prevents access to protected routes
4. **Automatic Cleanup**: Tokens cleared on logout
5. **Error Boundaries**: Graceful error handling

## Usage Examples

### Protecting a Route
```typescript
import { useRequireAuth } from '@/hooks/useAuth'

export default function MyPage() {
  const { isReady, user } = useRequireAuth()
  
  if (!isReady) return <Loading />
  
  return <div>Hello {user.name}!</div>
}
```

### Making Authenticated API Calls
```typescript
// Token is automatically added by interceptor
const response = await api.get('/protected-endpoint')
```

### Accessing User Data
```typescript
const { user, isAuthenticated } = useAuthStore()

if (isAuthenticated) {
  console.log(`Welcome ${user.name}!`)
}
```

## Testing

Run the authentication test script:
```bash
chmod +x scripts/test_auth.sh
./scripts/test_auth.sh
```

This will test:
- User registration
- User login
- Protected endpoint access
- Token persistence
- Error handling

## Environment Variables

Create `web/.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME="AI Analytics Platform"
```

## Demo Credentials

For testing purposes, you can use:
- **Email**: demo@example.com
- **Password**: demo123

Or create a new account through the registration form.
