# 7Solution Backend Challenge 🚀

## About This Project 📋

This is a user management API service built with Go. It provides basic user operations (create, read, update, delete) and connects to MongoDB for data storage.

## Features ✨

- 👤 **User Management**: Create, get, update and delete users
- 🔐 **JWT Authentication**: Secure API endpoints with JSON Web Tokens
- 📊 **MongoDB Database**: Store user data in MongoDB
- 🧮 **Concurrent User Counting**: Background goroutine logs total user count every 10 seconds
- 🐳 **Docker Support**: Run everything in containers for easy setup

## How to Run the Project 🏃‍♂️

### Requirements

- Docker and Docker Compose
- Go (for local development)
- MongoDB (already set up in Docker)
- .env file for local development
- .env.docker file for docker development

### Using the Makefile (Windows) 🪟

This project includes a Makefile designed for Windows users:

```bash
# Build the API binary
make build_api

# Start all Docker containers with fresh build
make up_build

# Start Docker containers without rebuilding
make up

# Stop all Docker containers
make down

# Run tests
make test

# Start API locally (not in Docker)
make start_api
```

### Alternative Commands for Mac/Linux Users 🍎

If you're using Mac or Linux, the Windows Makefile may not work directly. You can use these commands instead:

```bash
# Build the API binary
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o apiApp ./cmd/server

# Start Docker containers with fresh build
docker-compose up --build -d

# Start Docker containers without rebuilding
docker-compose up -d

# Stop all Docker containers
docker-compose down

# Run tests
go test ./...

# Start API locally (not in Docker)
go run ./cmd/server/main.go
```

### Environment Variables example 🛠️

The environment variables are as follows:

```env
APP_HOST=127.0.0.1 # 0.0.0.0 for docker
APP_PORT=3000
APP_NAME=your_app_name
APP_VERSION=v.0.1.0
APP_BODY_LIMIT=10490000 # max body size in bytes
APP_READ_TIMEOUT=60 # max read timeout in seconds
APP_WRITE_TIMEOUT=60 # max write timeout in seconds

JWT_SECRET_KEY=your_jwt_secret_key # jwt secret key
JWT_ACCESS_EXPIRES=86400 # jwt access expires in seconds

DB_HOST=db
DB_PORT=27017
DB_NAME=userdb
DB_USER=root
DB_PASSWORD=root
DB_MAX_POOL_SIZE=25
```


## API Endpoints 🌐

When the application is running, you can access these endpoints:

- `GET /v1/users`: Get all users
- `GET /v1/users/:id`: Get a specific user
- `POST /v1/users`: Create a new user
- `PUT /v1/users/:id`: Update a user (Protected Endpoint)
- `DELETE /v1/users/:id`: Delete a user (Protected Endpoint)
- `POST /v1/users/login`: Login and get authentication token

## Project Structure 📚

```
7solution-be/
├── cmd/
│   └── server/
│       └── main.go          # Entry point of the application
├── internal/
│   ├── auth/
│   │   ├── jwt.go           # JWT implementation
│   │   └── middleware.go    # Authentication middleware
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   └── mongodb.go       # MongoDB connection and operations
│   ├── servers/
│   │   └── server.go        # API server setup
│   └── users/
│       ├── controller.go    # HTTP handlers for user endpoints
│       ├── model.go         # User data models
│       ├── repository.go    # Database operations for users
│       ├── route.go         # API route definitions
│       └── test/            # Test files for user module
│           └── repository_test.go
├── pkg/
│   └── response/          # Standardized API response utilities
├── .dockerignore            # Files to exclude from Docker builds
├── .env                     # Local environment configuration
├── .env.docker              # Docker environment configuration
├── Dockerfile               # Instructions for building the API image
├── docker-compose.yml       # Multi-container Docker setup
└── Makefile                 # Build automation for the project
```

- **cmd**: Contains application entry points
- **internal**: Houses all internal application code
  - **auth**: Authentication and authorization logic
  - **config**: Application configuration handling
  - **database**: Database connections and common operations
  - **servers**: HTTP server setup and configuration
  - **users**: Complete user module with controller, model, and repository

## Environment Configuration ⚙️

The project uses two environment files:
- `.env`: For local development
- `.env.docker`: For Docker deployment

## Docker Setup 🐳

The project includes:
- `Dockerfile`: Builds the API container
- `docker-compose.yml`: Sets up both API and MongoDB services

When running in Docker, the API will be available at: http://localhost:3000

## Running Tests 🧪

To run all tests:
```bash
make test       # Windows
go test ./...   # Mac/Linux
```

## Development with Hot Reload 🔥

For local development, you can use Air for hot reloading:

```bash
# Install Air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run the API with hot reload
air
```

This will automatically rebuild and restart the API whenever you make code changes.

## Sample API Requests/Responses 📝

### Create User

**Request:**
```json
POST /v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "id": "683fc5ad99d4a0517a6a569f",
  "name": "John Doe",
  "email": "john@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI3c29sdXRpb24tYmUiLCJleHAiOjE3NDkwOTYyMzcsImlhdCI6MTc0OTAwOTgzNywiaXNzIjoiN3NvbHV0aW9uLWJlIiwibmJmIjoxNzQ5MDA5ODM3LCJzdWIiOiI2ODNmYzVhZDk5ZDRhMDUxN2E2YTU2OWYifQ.Ys8uAqmVI630l1OdDxyl02cfdeU-6yDRmHeGB5Wagh8"
}
```

### Login

**Request:**
```json
POST /v1/users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "id": "683fb6e9d1cf36b6bd03f6e4",
  "name": "John Doe",
  "email": "john@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI3c29sdXRpb24tYmUiLCJleHAiOjE3NDkwOTI1MTIsImlhdCI6MTc0OTAwNjExMiwiaXNzIjoiN3NvbHV0aW9uLWJlIiwibmJmIjoxNzQ5MDA2MTEyLCJzdWIiOiI2ODNmYjZlOWQxY2YzNmI2YmQwM2Y2ZTQifQ.-3y50PKTriBIU03jNLVYZ11_6x_TJNgnqv4tfsBRmJE"
}
```

### Get All Users

**Request:**
```
GET /v1/users
```

**Response:**
```json
[
 {
    "id": "683ef6567713989f89d4231c",
    "name": "Alex Johnson",
    "email": "alex@example.com"
  },
  {
    "id": "683ef6667713989f89d42320",
    "name": "Jane Smith",
    "email": "jane@example.com"
  }
]
```

### Get User By Id

**Request:**
```
GET /v1/users/683ef6567713989f89d4231c
```

**Response:**
```json
{
    "id": "683ef6567713989f89d4231c",
    "name": "Alex Johnson",
    "email": "alex@example.com"
}
```

### Update User

**Request:**
```json
PUT /v1/users/507f1f77bcf86cd799439011
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "name": "John Wilson",
  "email": "johnwilson@example.com"
}
```

**Response:**
```json
{
    "id": "507f1f77bcf86cd799439011",
    "name": "John Wilson",
    "email": "johnwilson@example.com",
    "message": "User updated successfully"
}
```

### Delete User

**Request:**
```
DELETE /v1/users/507f1f77bcf86cd799439011
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```text
User with id 507f1f77bcf86cd799439011 deleted successfully
```

## Assumptions and Design Decisions 🧰

1. **PUT vs PATCH**: The update user endpoint uses PUT instead of PATCH and requires all fields (name and email) to be provided, as it's a complete replacement of the resource rather than a partial update.

2. **JWT Authentication Method**: This api assumes that the required authentication method is jwt header-based authentication (`Authorization: Bearer token`) rather than cookies.

3. **Email Uniqueness Check**: The email field should be unique in the database. but assuming the database doesn't have the constraint, the application will handle the uniqueness check.

## Troubleshooting 🔧

If you see connection errors to MongoDB:
- Check if MongoDB container is running (`docker ps`)
- Verify environment settings in `.env.docker`
- Make sure `DB_HOST=db` in the Docker environment

## Future Improvements 🚀

- ⚡ **Rate Limiting**: Add protection against excessive API requests
- 🌍 **CORS Handling**: Improved cross-origin resource sharing for web clients
- 🔍 **Request ID Generation**: Unique IDs for each request to improve logging and debugging
- 📦 **Enhanced Docker Setup**: Use Docker secrets instead of copying env files into containers
- 📄 **Pagination & Total Count**: Enhance GET /users endpoint with pagination parameters and total count in response for better client-side handling
- 🛡️ **Role Check**: Add role check for protected endpoints (eg. only admin can delete users)

## Thank you for your consideration 🙏
