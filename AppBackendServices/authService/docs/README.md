# Authentication Service API Documentation

This directory contains the OpenAPI specification for the Authentication Service API.

## Overview

The Authentication Service provides endpoints for user authentication, registration, and account management. It supports:

- User registration and login
- JWT-based authentication
- Password management (reset, change)
- Biometric authentication
- User profile management

## OpenAPI Specification

The API is documented using the OpenAPI 3.0.3 specification in `openapi.yaml`. This specification can be used to:

- Generate client libraries
- Create API documentation
- Test API endpoints
- Validate API requests and responses

## Using the OpenAPI Documentation

### Viewing the Documentation

You can view the API documentation using tools like:

- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [ReDoc](https://redocly.github.io/redoc/)
- [Stoplight Studio](https://stoplight.io/studio)

To use Swagger UI, you can run:

```bash
# Install swagger-ui-cli
npm install -g swagger-ui-cli

# Serve the documentation
swagger-ui-cli serve docs/openapi.yaml
```

### Generating Client Libraries

You can generate client libraries for various programming languages using tools like:

- [OpenAPI Generator](https://openapi-generator.tech/)
- [Swagger Codegen](https://swagger.io/tools/swagger-codegen/)

Example using OpenAPI Generator:

```bash
# Install OpenAPI Generator
npm install @openapitools/openapi-generator-cli -g

# Generate a Go client
openapi-generator-cli generate -i docs/openapi.yaml -g go -o ./clients/go

# Generate a JavaScript client
openapi-generator-cli generate -i docs/openapi.yaml -g javascript -o ./clients/js
```

### Testing the API

You can use tools like [Postman](https://www.postman.com/) or [Insomnia](https://insomnia.rest/) to test the API endpoints. Both tools can import the OpenAPI specification to automatically create a collection of requests.

## API Endpoints

The API provides the following endpoints:

### Authentication

- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - User logout

### Password Management

- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password
- `POST /auth/change-password` - Change password

### Biometric Authentication

- `POST /auth/biometric/register` - Register biometric data
- `POST /auth/biometric/verify` - Verify biometric data

### User Profile

- `GET /auth/profile` - Get user profile
- `PUT /auth/profile` - Update user profile

### Health Check

- `GET /health` - Check service health

## Authentication

Most endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token:

```
Authorization: Bearer <token>
```

## Error Handling

The API uses standard HTTP status codes and returns error responses in the following format:

```json
{
  "code": 400,
  "message": "Invalid request",
  "details": {
    "field": "email",
    "reason": "Invalid email format"
  }
}
```

## Versioning

The API is versioned using the URL path. The current version is v1, which is implicit in the paths.

## Rate Limiting

The API implements rate limiting to prevent abuse. Rate limits are communicated in the response headers:

- `X-RateLimit-Limit` - Maximum number of requests per time window
- `X-RateLimit-Remaining` - Number of requests remaining in the current time window
- `X-RateLimit-Reset` - Time when the rate limit resets (Unix timestamp)

## Support

For API support, please contact support@example.com. 