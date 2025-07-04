# Car Rental API

A comprehensive Go-based car rental management system API built using Gin framework with PostgreSQL database. Features advanced booking management, membership system with discounts, driver services, and automatic stock management.

## Documentation

> **⚡ Quick Access**: **[Complete API Documentation](docs/api-documentation.md)** 
> 
> Comprehensive documentation with examples, validation rules, business logic, and all endpoint details.

## 📑 Table of Contents

- [🚀 Quick Start](#🚀-quick-start)
- [📋 API Endpoints](#📋-api-endpoints)
- [✨ Key Features](#✨-key-features)
- [🗄️ Database Schema](#️🗄️-database-schema)
- [⚙️ Technology Stack](#️⚙️-technology-stack)
- [📁 Project Structure](#📁-project-structure)

## 🚀 Quick Start

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

## 📋 API Endpoints

### Health Check
- `GET /health` - API status check

### Customer Management
- `GET /customers` - List all customers with membership details
- `GET /customers/:id` - Get customer by ID with membership information
- `POST /customers` - Create new customer
- `PUT /customers/:id` - Update customer information
- `DELETE /customers/:id` - Delete customer
- `PUT /customers/:id/subscribe/:membership_id` - Subscribe customer to membership
- `DELETE /customers/:id/unsubscribe` - Remove customer membership

### Car Management
- `GET /cars` - List all cars
- `GET /cars/:id` - Get car by ID
- `POST /cars` - Create new car
- `PUT /cars/:id` - Update car
- `DELETE /cars/:id` - Delete car

### Booking Management
- `GET /bookings` - List all bookings with complete details (customer, car, driver, booking type)
- `GET /bookings/:id` - Get booking by ID
- `POST /bookings` - Create new booking
- `PUT /bookings/:id` - Update booking
- `DELETE /bookings/:id` - Delete booking (restores car stock)
- `PUT /bookings/:id/finish` - Mark booking as finished
- `GET /bookings/types` - List all booking types (Car Only, Car & Driver)
- `GET /bookings/types/:id` - Get specific booking type

### Membership Management (Read-Only)
- `GET /memberships` - List all available memberships
- `GET /memberships/:id` - Get membership details with discount information

### Driver Management
- `GET /drivers` - List all drivers with availability status
- `GET /drivers/:id` - Get driver by ID
- `POST /drivers` - Create new driver
- `PUT /drivers/:id` - Update driver
- `DELETE /drivers/:id` - Delete driver
- `GET /drivers/:id/incentives` - Get driver incentive history

## ✨ Key Features

- **🏷️ Membership System**: Customer membership with automatic discount calculation
- **🚗 Driver Services**: Car & Driver booking options with driver cost calculation  
- **📦 Automatic Stock Management**: Car inventory automatically updated on booking operations
- **💰 Cost Calculation**: Total costs computed with rental duration and daily rates
- **🔄 Booking Types**: Support for "Car Only" and "Car & Driver" booking types
- **✅ Advanced Validation**: Constraint-based validation with detailed error responses
- **🔗 Relationship Management**: Comprehensive foreign key handling and referential integrity

<details>
<summary><strong>📋 Detailed Features & Business Logic</strong></summary>

### Advanced Features

**Stock Management**
- Car inventory automatically decremented on booking creation
- Stock restored on booking deletion or completion
- Prevents overbooking with availability checking

**Cost Calculation**
- Base cost: (rental days) × (car daily rent)
- Membership discounts automatically applied
- Driver costs calculated and added for Car & Driver bookings
- Total cost includes all applicable fees and discounts

**Validation & Constraints**
- Customer and car existence validation
- Car availability checking (stock > 0)
- Date validation (start date cannot be in past, must be before end date)
- Booking modification restrictions (cannot modify finished bookings)
- NIK uniqueness and format validation (16 characters)
- Phone number format validation (max 15 characters)

**Membership Integration**
- Customers can subscribe/unsubscribe to memberships
- Automatic discount application during booking cost calculation
- Membership details included in customer and booking responses

**Driver Assignment**
- Optional driver assignment for Car & Driver bookings
- Driver incentive history tracking
- Separate cost calculation for driver services

</details>

## 🗄️ Database Schema

![erd-v2.jpeg](docs/erd-v2.jpeg)

<details>
<summary><strong>🗄️ Database Structure</strong></summary>

### Customer Table
- **no** (PK) - `int` - Primary key, unique customer identifier
- **name** - `varchar` - Customer's full name (required)
- **nik** - `varchar(16)` - National identification number (required, unique, 16 chars)
- **phone_number** - `varchar(15)` - Customer's contact phone number (required, max 15 chars)
- **membership_id** (FK) - `int` - Foreign key referencing Membership.no (optional)

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
- **discount** - `float` - Applied discount amount (default: 0)
- **booking_type_id** (FK) - `int` - Foreign key referencing BookingType.no (required)
- **driver_id** (FK) - `int` - Foreign key referencing Driver.no (optional)
- **total_driver_cost** - `float` - Total driver cost (default: 0)

### Membership Table
- **no** (PK) - `int` - Primary key, unique membership identifier
- **membership_name** - `varchar` - Name of the membership tier (required)
- **discount** - `float` - Discount percentage offered by membership (required)

### Driver Table
- **no** (PK) - `int` - Primary key, unique driver identifier
- **name** - `varchar` - Driver's full name (required)
- **nik** - `varchar(16)` - National identification number (required, unique)
- **phone_number** - `varchar(15)` - Driver's contact phone number (required)
- **license_number** - `varchar` - Driver's license number (required, unique)
- **daily_rate** - `float` - Daily rate for driver services (required)
- **available** - `bool` - Driver availability status (default: true)

### BookingType Table
- **no** (PK) - `int` - Primary key, unique booking type identifier
- **name** - `varchar` - Type name (e.g., "Car Only", "Car & Driver")
- **description** - `varchar` - Description of the booking type

### DriverIncentive Table
- **no** (PK) - `int` - Primary key, unique incentive identifier
- **driver_id** (FK) - `int` - Foreign key referencing Driver.no (required)
- **amount** - `float` - Incentive amount (required)
- **date** - `datetime` - Date when incentive was awarded (required)
- **description** - `varchar` - Description of the incentive reason

### Relationships
1. **Customer → Booking**: One-to-Many (A customer can have multiple bookings)
2. **Customer → Membership**: Many-to-One (Multiple customers can have the same membership)
3. **Cars → Booking**: One-to-Many (A car model can be booked multiple times)
4. **Driver → Booking**: One-to-Many (A driver can be assigned to multiple bookings)
5. **Driver → DriverIncentive**: One-to-Many (A driver can have multiple incentives)
6. **BookingType → Booking**: One-to-Many (A booking type can be used for multiple bookings)

</details>

## ⚙️ Technology Stack

- **Framework**: [Gin](https://gin-gonic.com/) - HTTP web framework
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Environment**: [godotenv](https://github.com/joho/godotenv) for configuration
- **Language**: Go 1.23.1

## 📁 Project Structure

```
car-rental-v2/
├── cmd/
│   └── main.go              # Application entry point
├── pkg/
│   ├── database/            # Database connection and seeding
│   │   ├── database.go      # Database configuration and connection
│   │   └── seed.go          # Database seeding with initial data
│   ├── handlers/            # HTTP request handlers
│   │   ├── customer.go      # Customer CRUD + membership operations
│   │   ├── car.go          # Car CRUD operations
│   │   ├── booking.go      # Booking CRUD + finish operations
│   │   ├── membership.go   # Membership read operations
│   │   ├── driver.go       # Driver CRUD + incentive operations
│   │   └── booking_type.go # Booking type read operations
│   ├── models/              # Data models and validation
│   │   ├── customer.go     # Customer model
│   │   ├── car.go         # Car model
│   │   ├── booking.go     # Booking model
│   │   ├── membership.go  # Membership model
│   │   ├── driver.go      # Driver model
│   │   ├── driver_incentive.go # Driver incentive model
│   │   └── booking_type.go # Booking type model
│   ├── routes/              # API route definitions
│   │   └── routes.go
│   └── utils/               # Utility functions
│       └── referential_integrity.go # Database constraint utilities
├── docs/
│   ├── erd-v1.jpeg         # Original Entity Relationship Diagram
│   ├── erd-v2.jpeg         # Updated Entity Relationship Diagram
│   └── api-documentation.md # Complete API documentation
├── bin/
│   └── car-rental.exe      # Compiled binary
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
└── README.md
```