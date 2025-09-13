# User Management REST API Design

## Overview
A simple REST API for managing users with CRUD operations, following RESTful principles.

## Base URL
```
https://api.example.com/v1/
```

## Endpoints

### Users Resource
```
GET    /users              # List all users (paginated)
POST   /users              # Create new user
GET    /users/{id}         # Get specific user
PUT    /users/{id}         # Update entire user
PATCH  /users/{id}         # Partial update user
DELETE /users/{id}         # Delete user
```

## User Data Model
```json
{
  "id": "string (UUID)",
  "name": "string",
  "email": "string (unique)",
  "role": "string (optional: admin, user, etc.)",
  "created_at": "string (ISO 8601)",
  "updated_at": "string (ISO 8601)"
}
```

## Request/Response Examples

### Create User (POST /users)
**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user"
}
```

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### List Users (GET /users)
**Response (200 OK):**
```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 1
  }
}
```

## Authentication
- Use JWT tokens in Authorization header
- Required for all endpoints except GET /users (public read)

## Error Handling
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing/invalid token
- 403 Forbidden: Insufficient permissions
- 404 Not Found: User not found
- 409 Conflict: Email already exists

## Validation Rules
- Name: Required, 2-100 characters
- Email: Required, valid email format, unique
- Role: Optional, enum: admin, user, moderator

## Rate Limiting
- 100 requests per hour per user
- 1000 requests per hour for admin users

This design provides a clean, scalable foundation for user management with room for future enhancements.