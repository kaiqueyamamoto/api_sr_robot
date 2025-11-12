# JWT Authentication API Examples

## Overview

This API provides JWT-based authentication with email and password. Tokens expire after 24 hours (1 day).

## Base URL

```
http://localhost:8080
```

## Endpoints

### 0. Health Check

**GET** `/health`

Check if the server is healthy and running.

**Response (200 OK):**

```json
{
  "status": "healthy",
  "database": "connected"
}
```

**cURL Example:**

```bash
curl http://localhost:8080/health
```

---

### Metrics Endpoint

**GET** `/metrics`

Prometheus metrics endpoint for monitoring (see PROMETHEUS_METRICS.md for details).

**cURL Example:**

```bash
curl http://localhost:8080/metrics
```

---

### 1. Register a New User

**POST** `/auth/register`

Register a new user with email and password.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T17:10:57.738Z"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

---

### 2. Login

**POST** `/auth/login`

Login with existing credentials to get a JWT token.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T17:10:57.738Z"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

---

### 3. User Info

**GET** `/userinfo`

Get information about the currently authenticated user including timestamps.

**Headers:**

```
Authorization: Bearer <your-jwt-token>
```

**Response (200 OK):**

```json
{
  "user_id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "created_at": "2025-11-12T17:10:57.738Z",
  "updated_at": "2025-11-12T17:10:57.738Z"
}
```

**cURL Example:**

```bash
curl http://localhost:8080/userinfo \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 4. Protected Endpoint (Chat)

**POST** `/chat`

This endpoint requires authentication. Include the JWT token in the Authorization header.

**Headers:**

```
Authorization: Bearer <your-jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Chat endpoint",
  "user_email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/chat \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid input data"
}
```

### 401 Unauthorized

```json
{
  "error": "Invalid or expired token"
}
```

### 409 Conflict

```json
{
  "error": "User already exists"
}
```

---

## Token Information

- **Expiration:** 1 day (24 hours) from issuance
- **Algorithm:** HS256
- **Format:** Bearer token in Authorization header

---

## Environment Variables

- `MONGODB_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)

---

## Running the Server

1. Make sure MongoDB is running
2. Run the server:

```bash
go run main.go
```

Or build and run:

```bash
go build -o chatserver main.go
./chatserver
```

---

## Full Workflow Example

```bash
# 1. Register a new user
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123456"}' \
  | jq -r '.token')

echo "Registered! Token: $TOKEN"

# 2. Use the token to access protected endpoint
curl -X POST http://localhost:8080/chat \
  -H "Authorization: Bearer $TOKEN"

# 3. Login again (if token expired or you need a new session)
NEW_TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123456"}' \
  | jq -r '.token')

echo "Logged in! New Token: $NEW_TOKEN"

# 4. Use the new token
curl -X POST http://localhost:8080/chat \
  -H "Authorization: Bearer $NEW_TOKEN"
```

---

## Security Notes

⚠️ **Important for Production:**

1. Change the JWT secret key in `controllers/auth.go` to a strong, random value
2. Store the JWT secret in an environment variable, not in the code
3. Use HTTPS in production
4. Consider adding rate limiting
5. Implement password strength requirements
6. Add email verification
7. Consider adding refresh tokens for longer sessions
