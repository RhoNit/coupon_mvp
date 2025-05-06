# Coupon System MVP

A robust coupon management system built with Go and Echo framework, designed for a medicine ordering platform.

## Features

- Admin Coupon Creation and Management
- Coupon Validation with Multiple Constraints
- Support for Different Usage Types (One-time, Multi-use, Time-based)
- Concurrent-safe Operations
- Redis Caching Implementation
- OpenAPI Documentation

## Architecture

The system follows a clean architecture pattern with the following layers:

- **API Layer**: Handles HTTP requests and responses using Echo framework
- **Service Layer**: Contains business logic for coupon management
- **Repository Layer**: Manages data persistence using PostgreSQL
- **Domain Layer**: Contains core business entities and interfaces

### Concurrency & Caching

- Uses mutex locks for concurrent coupon validations
- Implements Redis caching for frequently accessed coupons
- Database-level safety with proper transaction handling

## Setup Instructions

1. **Prerequisites**
   - Go 1.21 or higher
   - PostgreSQL
   - Redis
   - Docker

2. **Installation**
   ```bash
   # Clone the repository
   git clone https://github.com/RhoNit/coupon_mvp.git
   cd coupon_mvp

   # Install dependencies
   go mod download

   # Set up environment variables
   cp .env.example .env
   # Edit .env with your configuration

   # Generate Swagger documentation
   swag init -g main.go -o docs

   # Run the application
   go run main.go
   ```

3. **Docker Setup**
   ```bash
   docker-compose up -d
   ```

## API Documentation

Swagger documentation is available at `http://localhost:8081/swagger/index.html` when the server is running.

### API Endpoints

#### 1. Create Coupon
```bash
curl -X POST http://localhost:8081/api/v1/coupons \
  -H "Content-Type: application/json" \
  -d '{
    "coupon_code": "SAVE20",
    "expiry_date": "2024-12-31T23:59:59Z",
    "usage_type": "multi_use",
    "applicable_medicine_ids": ["med_123", "med_456"],
    "applicable_categories": ["painkiller"],
    "min_order_value": 100.00,
    "valid_time_window": {
      "start_time": "2024-01-01T00:00:00Z",
      "end_time": "2024-12-31T23:59:59Z"
    },
    "terms_and_conditions": "Valid on all painkillers",
    "discount_type": "percentage",
    "discount_value": 20.00,
    "max_usage_per_user": 3
  }'
```

**Response:**
```json
{
  "id": "uuid-generated",
  "coupon_code": "SAVE20",
  "expiry_date": "2024-12-31T23:59:59Z",
  "usage_type": "multi_use",
  "applicable_medicine_ids": ["med_123", "med_456"],
  "applicable_categories": ["painkiller"],
  "min_order_value": 100.00,
  "valid_time_window": {
    "start_time": "2024-01-01T00:00:00Z",
    "end_time": "2024-12-31T23:59:59Z"
  },
  "terms_and_conditions": "Valid on all painkillers",
  "discount_type": "percentage",
  "discount_value": 20.00,
  "max_usage_per_user": 3,
  "created_at": "2024-03-15T10:00:00Z",
  "updated_at": "2024-03-15T10:00:00Z"
}
```

#### 2. Get Applicable Coupons
```bash
curl -X GET http://localhost:8081/api/v1/coupons/applicable \
  -H "Content-Type: application/json" \
  -d '{
    "cart_items": [
      {
        "id": "med_123",
        "category": "painkiller",
        "price": 150.00
      }
    ],
    "order_total": 150.00,
    "timestamp": "2024-03-15T10:00:00Z"
  }'
```

**Response:**
```json
[
  {
    "coupon_code": "SAVE20",
    "discount_value": 20.00
  }
]
```

#### 3. Validate Coupon
```bash
curl -X POST http://localhost:8081/api/v1/coupons/validate \
  -H "Content-Type: application/json" \
  -d '{
    "coupon_code": "SAVE20",
    "cart_items": [
      {
        "id": "med_123",
        "category": "painkiller",
        "price": 150.00
      }
    ],
    "order_total": 150.00,
    "timestamp": "2024-03-15T10:00:00Z"
  }'
```

**Success Response:**
```json
{
  "is_valid": true,
  "discount": {
    "items_discount": 30.00,
    "charges_discount": 0.00
  },
  "message": "Coupon applied successfully"
}
```

**Error Response:**
```json
{
  "is_valid": false,
  "reason": "coupon expired or not applicable"
}
```

### Error Codes

- `400 Bad Request`: Invalid request format
- `404 Not Found`: Coupon not found
- `409 Conflict`: Coupon already exists
- `422 Unprocessable Entity`: Validation error
- `500 Internal Server Error`: Server error

### Data Models

#### Coupon
```json
{
  "id": "string",
  "coupon_code": "string",
  "expiry_date": "datetime",
  "usage_type": "one_time|multi_use|time_based",
  "applicable_medicine_ids": ["string"],
  "applicable_categories": ["string"],
  "min_order_value": "decimal",
  "valid_time_window": {
    "start_time": "datetime",
    "end_time": "datetime"
  },
  "terms_and_conditions": "string",
  "discount_type": "percentage|fixed",
  "discount_value": "decimal",
  "max_usage_per_user": "integer",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

#### CartItem
```json
{
  "id": "string",
  "category": "string",
  "price": "decimal"
}
```

## Database Schema

The system uses PostgreSQL with the following main tables:
- coupons
- coupon_usage
- medicines
- categories

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./internal/service -v
```

## License

MIT 