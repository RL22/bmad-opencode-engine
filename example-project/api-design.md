# REST API Design Structure

## Overview
This document outlines a simple, scalable REST API structure following industry best practices and RESTful principles.

## Base URL Structure
```
https://api.example.com/v1/
```

## Core Principles

### 1. Resource-Based URLs
- Use nouns, not verbs in URLs
- Use plural nouns for collections
- Use hierarchical structure for relationships

### 2. HTTP Methods
- `GET` - Retrieve data (safe, idempotent)
- `POST` - Create new resources
- `PUT` - Update/replace entire resource (idempotent)
- `PATCH` - Partial update of resource
- `DELETE` - Remove resource (idempotent)

### 3. Status Codes
- `200` - OK (successful GET, PUT, PATCH)
- `201` - Created (successful POST)
- `204` - No Content (successful DELETE)
- `400` - Bad Request (client error)
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Unprocessable Entity (validation errors)
- `500` - Internal Server Error

## API Endpoints Structure

### Users Resource
```
GET    /users              # List all users (with pagination)
POST   /users              # Create new user
GET    /users/{id}         # Get specific user
PUT    /users/{id}         # Update entire user
PATCH  /users/{id}         # Partial update user
DELETE /users/{id}         # Delete user
```

### Posts Resource (Example nested resource)
```
GET    /posts              # List all posts
POST   /posts              # Create new post
GET    /posts/{id}         # Get specific post
PUT    /posts/{id}         # Update entire post
PATCH  /posts/{id}         # Partial update post
DELETE /posts/{id}         # Delete post

# User's posts (nested resource)
GET    /users/{id}/posts   # Get posts by specific user
POST   /users/{id}/posts   # Create post for specific user
```

### Comments Resource (Nested under posts)
```
GET    /posts/{id}/comments     # Get comments for a post
POST   /posts/{id}/comments     # Create comment on a post
GET    /comments/{id}           # Get specific comment
PUT    /comments/{id}           # Update comment
DELETE /comments/{id}           # Delete comment
```

## Request/Response Format

### Content Type
- Request: `Content-Type: application/json`
- Response: `Content-Type: application/json`

### Request Structure
```json
{
  "data": {
    "type": "user",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### Response Structure
```json
{
  "data": {
    "id": "123",
    "type": "user",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  },
  "meta": {
    "timestamp": "2025-01-01T00:00:00Z"
  }
}
```

### Collection Response (with pagination)
```json
{
  "data": [
    {
      "id": "123",
      "type": "user",
      "attributes": {
        "name": "John Doe",
        "email": "john@example.com"
      }
    }
  ],
  "meta": {
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total_pages": 5,
      "total_count": 100
    },
    "timestamp": "2025-01-01T00:00:00Z"
  },
  "links": {
    "self": "/users?page=1",
    "next": "/users?page=2",
    "prev": null,
    "first": "/users?page=1",
    "last": "/users?page=5"
  }
}
```

### Error Response
```json
{
  "errors": [
    {
      "id": "validation_error",
      "status": "422",
      "code": "INVALID_EMAIL",
      "title": "Invalid email format",
      "detail": "The email field must be a valid email address",
      "source": {
        "pointer": "/data/attributes/email"
      }
    }
  ],
  "meta": {
    "timestamp": "2025-01-01T00:00:00Z"
  }
}
```

## Query Parameters

### Filtering
```
GET /users?filter[status]=active
GET /users?filter[created_at][gte]=2025-01-01
```

### Sorting
```
GET /users?sort=name
GET /users?sort=-created_at  # Descending
GET /users?sort=name,-created_at  # Multiple fields
```

### Pagination
```
GET /users?page=2&per_page=50
```

### Field Selection (Sparse Fieldsets)
```
GET /users?fields[user]=name,email
```

### Including Related Resources
```
GET /posts?include=author,comments
```

## Authentication

### Bearer Token (Recommended)
```
Authorization: Bearer <token>
```

### API Key (Alternative)
```
X-API-Key: <api_key>
```

## Rate Limiting

### Headers
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

### Rate Limit Exceeded Response
```json
{
  "errors": [
    {
      "status": "429",
      "title": "Rate limit exceeded",
      "detail": "You have exceeded the rate limit of 1000 requests per hour"
    }
  ]
}
```

## Versioning

### URL Versioning (Recommended)
```
https://api.example.com/v1/users
https://api.example.com/v2/users
```

### Header Versioning (Alternative)
```
Accept: application/vnd.api+json;version=1
```

## CORS Headers
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Health Check Endpoint
```
GET /health

Response:
{
  "status": "healthy",
  "timestamp": "2025-01-01T00:00:00Z",
  "version": "1.0.0",
  "services": {
    "database": "healthy",
    "cache": "healthy"
  }
}
```

## Documentation Endpoint
```
GET /docs           # API documentation (Swagger/OpenAPI)
GET /schema         # JSON Schema for API
```

## Best Practices Summary

1. **Consistency**: Use consistent naming conventions and response formats
2. **Stateless**: Each request should contain all necessary information
3. **Cacheable**: Use appropriate HTTP headers for caching
4. **Layered**: Design with proxy servers and load balancers in mind
5. **Self-Descriptive**: Include metadata and links in responses
6. **HATEOAS**: Include links to related resources where appropriate

## Security Considerations

1. **HTTPS Only**: Never serve API over HTTP in production
2. **Input Validation**: Validate all input data
3. **Rate Limiting**: Implement rate limiting to prevent abuse
4. **Authentication**: Require authentication for sensitive operations
5. **Authorization**: Implement proper access controls
6. **Audit Logging**: Log all API access and modifications

This structure provides a solid foundation for building scalable, maintainable REST APIs that follow industry standards and best practices.