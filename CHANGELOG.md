# Changelog

## [Unreleased] - 2025-11-12

### Added

- **User timestamps**: Added `created_at` and `updated_at` fields to User model
- **User info endpoint**: New `GET /userinfo` endpoint to retrieve current user information
- Timestamps are automatically set during user registration
- Auth responses now include `user_id` and `created_at` fields

### Changed

- Updated `models/user.go` to include `CreatedAt` and `UpdatedAt` fields
- Updated `models/AuthResponse` to include `UserID` and `CreatedAt` fields
- Updated auth controller to set timestamps on user creation
- Updated registration and login responses to include user_id and created_at

### API Changes

#### POST /auth/register

Response now includes:

```json
{
  "token": "...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T17:10:57.738Z"
}
```

#### POST /auth/login

Response now includes:

```json
{
  "token": "...",
  "email": "user@example.com",
  "user_id": "507f1f77bcf86cd799439011",
  "created_at": "2025-11-12T17:10:57.738Z"
}
```

#### GET /userinfo (New)

Returns complete user information:

```json
{
  "user_id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "created_at": "2025-11-12T17:10:57.738Z",
  "updated_at": "2025-11-12T17:10:57.738Z"
}
```

## [1.0.0] - 2025-11-12

### Added

- JWT authentication with email and password
- 24-hour token expiration
- MongoDB integration
- Prometheus metrics
- Health check endpoint
- Authentication middleware
- Password hashing with bcrypt
- Comprehensive API documentation
- Docker Compose support for monitoring stack
- Grafana provisioning configuration
