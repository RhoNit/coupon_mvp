# Coupon System MVP

A robust coupon management system built with Go and Echo framework, designed for a medicine ordering platform.

## Features

- Admin Coupon Creation and Management
- Coupon Validation with Multiple Constraints
- Support for Different Usage Types (One-time, Multi-use, Time-based)
- Concurrent-safe Operations
- Caching Implementation
- OpenAPI Documentation

## Architecture

The system follows a clean architecture pattern with the following layers:

- **API Layer**: Handles HTTP requests and responses using Echo framework
- **Service Layer**: Contains business logic for coupon management
- **Repository Layer**: Manages data persistence using PostgreSQL
- **Domain Layer**: Contains core business entities and interfaces

### Concurrency & Caching

- Uses mutex locks for concurrent coupon validations
- Implements LRU caching for frequently accessed coupons
- Database-level safety with proper transaction handling

## Setup Instructions

1. **Prerequisites**
   - Go 1.21 or higher
   - PostgreSQL
   - Docker (optional)

2. **Installation**
   ```bash
   # Clone the repository
   git clone <repository-url>
   cd farmako_assignment

   # Install dependencies
   go mod download

   # Set up environment variables
   cp .env.example .env
   # Edit .env with your configuration

   # Run the application
   go run main.go
   ```

3. **Docker Setup**
   ```bash
   docker-compose up -d
   ```

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the server is running.

## Database Schema

The system uses PostgreSQL with the following main tables:
- coupons
- coupon_usage
- medicines
- categories

## Testing

```bash
go test ./...
```

## License

MIT 