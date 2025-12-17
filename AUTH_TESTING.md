# Authentication API Testing Guide

## Endpoints Added

Your Go API server now has three new authentication endpoints:

1. **POST /signup** - Register a new user
2. **POST /login** - Login with existing credentials
3. **POST /logout** - Logout (invalidate session)

---

## Testing the Endpoints

### 1. Sign Up a New User

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2025-12-16T10:30:00Z"
  }
}
```

**Save the token** - You'll need it for the logout endpoint!

---

### 2. Login with Existing User

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2025-12-16T10:30:00Z"
  }
}
```

---

### 3. Logout

```bash
# Replace YOUR_JWT_TOKEN with the actual token from signup/login
curl -X POST http://localhost:8080/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "message": "Successfully logged out",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Error Scenarios

### Signup with Existing Email
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "newpassword"
  }'
```

**Response (409 Conflict):**
```json
{
  "error": "User with this email already exists"
}
```

### Login with Wrong Password
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "wrongpassword"
  }'
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

### Logout without Token
```bash
curl -X POST http://localhost:8080/logout
```

**Response (401 Unauthorized):**
```json
{
  "error": "No authorization header provided"
}
```

---

## Architecture Overview

### Files Created:

1. **`internal/models/user.go`** - Data structures for User, SignupRequest, LoginRequest, AuthResponse
2. **`internal/database/memory.go`** - In-memory database using Go maps with thread-safe access
3. **`internal/utils/jwt.go`** - JWT token generation and validation functions
4. **`internal/handler/auth.go`** - HTTP handlers for signup, login, logout

### Key Concepts:

- **Password Hashing**: Uses bcrypt to securely hash passwords (never store plain text!)
- **JWT Authentication**: Generates JSON Web Tokens for stateless authentication
- **Thread Safety**: Uses `sync.RWMutex` for concurrent database access
- **Validation**: Gin's binding validates email format and password length automatically

---

## Next Steps

1. **Add Protected Routes**: Create middleware to verify JWT tokens on protected endpoints
2. **Refresh Tokens**: Implement refresh token mechanism for better security
3. **Real Database**: Replace in-memory DB with PostgreSQL, MySQL, or MongoDB
4. **Environment Variables**: Move JWT secret to environment variables
5. **Token Blacklist**: Implement Redis-based token blacklist for logout

---

## Security Notes (Production Checklist)

- ⚠️ **Change JWT Secret**: Use a long, random secret from environment variables
- ⚠️ **HTTPS Only**: Always use HTTPS in production to encrypt tokens in transit
- ⚠️ **Token Expiration**: Consider shorter expiration times (15 min) with refresh tokens
- ⚠️ **Rate Limiting**: Add rate limiting to prevent brute force attacks
- ⚠️ **Input Validation**: Add more robust email/password validation
- ⚠️ **CORS**: Configure CORS properly for your frontend domain
