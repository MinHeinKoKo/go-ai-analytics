# Authentication System

This project includes a complete JWT-based authentication system with user registration, login, and protected routes.

## Features

- User registration with email and password
- User login with JWT token generation
- Password hashing using bcrypt
- JWT middleware for protecting routes
- Input validation
- MongoDB integration with indexes

## API Endpoints

### Authentication Routes

#### Register User
```
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login User
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get Current User
```
GET /api/auth/me
Authorization: Bearer <token>
```

Response:
```json
{
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Protected Routes Examples

#### Get Profile
```
GET /api/protected/profile
Authorization: Bearer <token>
```

#### Get Dashboard
```
GET /api/protected/dashboard
Authorization: Bearer <token>
```

## Environment Variables

Make sure to set these environment variables in your `.env` file:

```env
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
MONGODB_URI=mongodb://localhost:27017
MONGODB_DBNAME=ai-analytics
```

## Usage in Frontend

### JavaScript/TypeScript Example

```javascript
// Register user
const registerUser = async (userData) => {
  const response = await fetch('/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('token', data.token);
    return data;
  }
  throw new Error(data.error);
};

// Login user
const loginUser = async (credentials) => {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('token', data.token);
    return data;
  }
  throw new Error(data.error);
};

// Make authenticated requests
const makeAuthenticatedRequest = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
};
```

## Security Features

- Passwords are hashed using bcrypt with default cost
- JWT tokens have configurable expiry time
- Email uniqueness is enforced at database level
- Input validation on all endpoints
- CORS configuration for frontend integration

## Database Schema

### Users Collection

```javascript
{
  _id: ObjectId,
  email: String (unique, required),
  password: String (hashed, required),
  first_name: String (required),
  last_name: String (required),
  is_active: Boolean (default: true),
  created_at: Date,
  updated_at: Date
}
```

## Testing

Run the authentication tests:

```bash
go test ./test/...
```

## Adding New Protected Routes

To create new protected routes, use the auth middleware:

```go
import "ai-analytics/internal/middleware"

// In your route handler
protected := r.Group("/api/your-protected-routes")
protected.Use(middleware.AuthMiddleware(config))
{
    protected.GET("/example", yourHandler)
}

// In your handler function
func yourHandler(c *gin.Context) {
    userID := c.MustGet("user_id").(primitive.ObjectID)
    userEmail := c.MustGet("user_email").(string)
    
    // Your protected logic here
}
```
