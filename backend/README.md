# Jobsity Backend API

A REST API server built with Go and Fiber framework.

## Features

- RESTful API endpoints
- CORS support
- Request logging
- Health check endpoint
- Login authentication with MongoDB
- User management

## Getting Started

### Prerequisites

- Go 1.19 or higher
- MongoDB (running on localhost:27017)
- Git

### Installation

1. Clone the repository
2. Navigate to the backend directory:

   ```bash
   cd backend
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Start MongoDB (if not already running):

   ```bash
   # On macOS with Homebrew
   brew services start mongodb-community

   # Or start manually
   mongod
   ```

5. Set up test users (optional):

   ```bash
   mongosh jobsity setup_test_user.js
   ```

6. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:3000`

## Project Structure

The project follows a clean architecture pattern with the following structure:

```
backend/
├── cmd/
│   └── server/          # Application entry point
├── pkg/
│   └── domain/          # Domain models and entities
│       └── user/        # User domain
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers
│   ├── repository/      # Repository layer (data access)
│   └── service/         # Service layer (business logic)
└── go.mod               # Go module file
```

This structure separates concerns:

- **pkg/domain**: Contains domain models and entities (public API)
- **internal/repository**: Data access layer (interfaces and implementations)
- **internal/service**: Business logic layer (interfaces and implementations)
- **internal/handlers**: HTTP request/response handling
- **internal/config**: Application configuration
- **internal/database**: Database connection management

## API Endpoints

### Health Check

- **GET** `/health`
- Returns server status

**Response:**

```json
{
  "status": "ok",
  "message": "Server is running"
}
```

### Login

- **POST** `/api/v1/login`
- Authenticates user credentials

**Request Body:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Login successful",
  "token": "fake-jwt-token-12345",
  "user": {
    "id": "user_id",
    "username": "admin",
    "email": "admin@jobsity.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Response (401):**

```json
{
  "success": false,
  "message": "Invalid username or password"
}
```

**Error Response (400):**

```json
{
  "success": false,
  "message": "Username and password are required"
}
```

### Create User

- **POST** `/api/v1/users`
- Creates a new user in the database

**Request Body:**

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "User created successfully",
  "user": {
    "id": "user_id",
    "username": "newuser",
    "email": "newuser@jobsity.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## Testing the API

### Health Check

```bash
curl -X GET http://localhost:3000/health
```

### Login with valid credentials

```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

### Login with invalid credentials

```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}'
```

### Create a new user

```bash
curl -X POST http://localhost:3000/create-user \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"newpass","email":"newuser@jobsity.com"}'
```

## Default Test Users

After running the setup script, you'll have these test users:

- Username: `admin`, Password: `password`
- Username: `testuser`, Password: `testpass`

**Note:** These are demo credentials. In a production environment, implement proper password hashing and user management.

## Testing

The project includes comprehensive tests for all layers:

### Running Tests

```bash
# Run all tests
make test-all

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage

# Or use the test script
./scripts/test.sh
```

### Test Structure

```
tests/
├── integration/          # Integration tests
│   └── api_test.go      # API endpoint tests
├── unit/                # Unit tests
│   ├── service_test.go  # Service layer tests
│   └── repository_test.go # Repository layer tests
├── test_config.go       # Test configuration
└── scripts/
    └── test.sh          # Test runner script
```

### Test Coverage

The tests cover:

- **API Endpoints**: All HTTP endpoints with various scenarios
- **Service Layer**: Business logic with mocked dependencies
- **Repository Layer**: Database operations with MongoDB mocks
- **Error Handling**: Invalid inputs and error scenarios
- **Authentication**: Login and user management flows

### Prerequisites for Testing

- MongoDB running on localhost:27017
- Go 1.19 or higher
- All dependencies installed (`make deps`)
