# Car Rental API

A Go-based car rental management system API built using Gin framework with PostgreSQL database.

## Quick Start

### Prerequisites
- Go 1.23.1 or higher
- PostgreSQL database

### Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
# Database Configuration
DB_HOST=localhost          # Default: localhost
DB_PORT=5432              # Default: 5432
DB_USER=postgres          # Default: postgres
DB_PASSWORD=your_password # Required
DB_NAME=car_rental        # Default: car_rental
DB_SSLMODE=disable        # Default: disable

# Server Configuration
PORT=8080                 # Default: 8080
GIN_MODE=debug           # Default: debug
```

4. Run the application:
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - API status check

### Customer Management
- `GET /customers` - List all customers
- `GET /customers/:id` - Get customer by ID
- `POST /customers` - Create new customer
- `PUT /customers/:id` - Update customer
- `DELETE /customers/:id` - Delete customer

### Car Management
- `GET /cars` - List all cars
- `GET /cars/:id` - Get car by ID
- `POST /cars` - Create new car
- `PUT /cars/:id` - Update car
- `DELETE /cars/:id` - Delete car

### Booking Management
- `GET /bookings` - List all bookings with customer and car details
- `GET /bookings/:id` - Get booking by ID
- `POST /bookings` - Create new booking
- `PUT /bookings/:id` - Update booking
- `DELETE /bookings/:id` - Delete booking
- `PUT /bookings/:id/finish` - Mark booking as finished

<details>
<summary><strong>📋 Detailed API Documentation</strong></summary>

### Complete API Documentation

**[📖 API Documentation](docs/api-documentation.md)** - Comprehensive documentation for all endpoints including:

- **Customer Management** - Full CRUD operations with validation rules
- **Car Management** - Inventory management with stock integration  
- **Booking Management** - Complete booking lifecycle with automatic cost calculation
- **Data Models** - Detailed field specifications and constraints
- **Error Handling** - Complete error scenarios and status codes
- **Business Rules** - Stock management, validation, and relationships

Key features include:
- **Automatic Stock Management**: Car inventory decremented on booking creation, restored on deletion/completion
- **Cost Calculation**: Total cost = (rental days) × (car daily rent)
- **Validation Rules**: 
  - Customer and car must exist
  - Car must be available (stock > 0)
  - Start date cannot be in the past
  - Start date must be before end date
  - Cannot modify/delete finished bookings
- **Error Handling**: Comprehensive error responses with appropriate HTTP status codes

</details>

## Database Schema

![erd-v1.jpeg](docs/erd-v1.jpeg)

<details>
<summary><strong>🗄️ Database Structure</strong></summary>

### Customer Table
- **no** (PK) - `int` - Primary key, unique customer identifier
- **name** - `varchar` - Customer's full name (required)
- **nik** - `varchar(16)` - National identification number (required, unique, 16 chars)
- **phone_number** - `varchar(15)` - Customer's contact phone number (required, max 15 chars)

### Cars Table
- **no** (PK) - `int` - Primary key, unique car identifier
- **name** - `varchar` - Car model/name (required)
- **stock** - `int` - Number of available cars of this model (required, min 0)
- **daily_rent** - `float` - Daily rental price (required, min 0)

### Booking Table
- **no** (PK) - `int` - Primary key, unique booking identifier
- **customer_id** (FK) - `int` - Foreign key referencing Customer.no (required)
- **cars_id** (FK) - `int` - Foreign key referencing Cars.no (required)
- **start_rent** - `datetime` - Rental start date and time (required)
- **end_rent** - `datetime` - Rental end date and time (required)
- **total_cost** - `float` - Total calculated cost for the rental period
- **finished** - `bool` - Flag indicating if the rental is completed (default: false)

### Relationships
1. **Customer → Booking**: One-to-Many (A customer can have multiple bookings)
2. **Cars → Booking**: One-to-Many (A car model can be booked multiple times)

</details>

## Technology Stack

- **Framework**: [Gin](https://gin-gonic.com/) - HTTP web framework
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Environment**: [godotenv](https://github.com/joho/godotenv) for configuration
- **Language**: Go 1.23.1

## Project Structure

```
car-rental-2/
├── cmd/
│   └── main.go              # Application entry point
├── pkg/
│   ├── database/            # Database connection and migrations
│   │   └── database.go      
│   ├── handlers/            # HTTP request handlers
│   │   ├── customer.go      # Customer CRUD operations
│   │   ├── car.go          # Car CRUD operations
│   │   └── booking.go      # Booking CRUD operations
│   ├── models/              # Data models and validation
│   │   ├── customer.go     # Customer model
│   │   ├── car.go         # Car model
│   │   └── booking.go     # Booking model
│   └── routes/              # API route definitions
│       └── routes.go       
├── docs/
│   ├── erd-v1.jpeg         # Entity Relationship Diagram
│   └── api-documentation.md # Complete API documentation
├── bin/
│   └── car-rental.exe      # Compiled binary
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
└── README.md              # This file
```