# Car Rental API Documentation

This document provides comprehensive documentation for all API endpoints in the car rental management system.

## Table of Contents

- [Overview](#overview)
- [Base URL](#base-url)
- [Health Check](#health-check)
- [Customer Endpoints](#customer-endpoints)
- [Car Endpoints](#car-endpoints)
- [Booking Endpoints](#booking-endpoints)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Business Rules](#business-rules)

## Overview

The Car Rental API provides complete CRUD operations for managing customers, cars, and bookings. The API features automatic stock management, cost calculation, and comprehensive validation.

### Key Features
- **Automatic Stock Management** - Car inventory automatically updated on booking operations
- **Cost Calculation** - Total costs computed based on rental duration and daily rates
- **Data Validation** - Input validation with detailed error messages
- **Relationship Management** - Proper foreign key handling and constraints

## Base URL

```
http://localhost:8080/api/v1
```

For health checks:
```
http://localhost:8080/health
```

---

## Health Check

### GET /health
Check API status and availability.

**Success Response (200 OK):**
```json
{
    "message": "Car Rental API",
    "status": "OK"
}
```

---

## Customer Endpoints

### GET /api/v1/customers
Retrieve all customers.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "name": "Wawan Hermawan",
            "nik": "3372093912739",
            "phone_number": "081237123682"
        },
        {
            "no": 2,
            "name": "Philip Walker",
            "nik": "3372093912785",
            "phone_number": "081237123683"
        },
        {
            "no": 3,
            "name": "Hugo Fleming",
            "nik": "3372093912800",
            "phone_number": "081237123684"
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve customers"
}
```

### GET /api/v1/customers/:id
Retrieve a specific customer by ID.

**URL Parameters:**
- `id` (integer) - Customer ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "name": "Wawan Hermawan",
        "nik": "3372093912739",
        "phone_number": "081237123682"
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid customer ID"
}

// 404 Not Found - Customer doesn't exist
{
    "error": "Customer not found"
}
```

### POST /api/v1/customers
Create a new customer.

**Request Body:**
```json
{
    "name": "John Doe",
    "nik": "1234567890123456",
    "phone_number": "081234567890"
}
```

**Field Requirements:**
- `name` (string, required) - Customer's full name
- `nik` (string, required) - National identification number (exactly 16 characters, must be unique)
- `phone_number` (string, required) - Customer's contact phone number (max 15 characters)

**Success Response (201 Created):**
```json
{
    "data": {
        "no": 4,
        "name": "John Doe",
        "nik": "1234567890123456",
        "phone_number": "081234567890"
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Missing required fields
{
    "error": "Key: 'Customer.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}

// 400 Bad Request - Invalid NIK length
{
    "error": "Key: 'Customer.NIK' Error:Field validation for 'NIK' failed on the 'len' tag"
}

// 500 Internal Server Error - Duplicate NIK
{
    "error": "Failed to create customer"
}
```

### PUT /api/v1/customers/:id
Update an existing customer.

**URL Parameters:**
- `id` (integer) - Customer ID

**Request Body:**
```json
{
  "name": "John Doe - Updated",
  "nik": "1234567890123456",
  "phone_number": "081234567890"
}
```

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 4,
        "name": "John Doe - Updated",
        "nik": "1234567890123456",
        "phone_number": "081234567890"
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid customer ID
{
    "error": "Invalid customer ID"
}

// 404 Not Found
{
    "error": "Customer not found"
}

// 500 Internal Server Error
{
    "error": "Failed to update customer"
}
```

### DELETE /api/v1/customers/:id
Delete a customer.

**URL Parameters:**
- `id` (integer) - Customer ID

**Success Response (200 OK):**
```json
{
    "message": "Customer deleted successfully"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid customer ID
{
    "error": "Invalid customer ID"
}

// 404 Not Found
{
    "error": "Customer not found"
}

// 500 Internal Server Error - Foreign key constraint
{
    "error": "Failed to delete customer"
}
```

---

## Car Endpoints

### GET /api/v1/cars
Retrieve all cars.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "name": "Toyota Camry",
            "stock": 2,
            "daily_rent": 500000
        },
        {
            "no": 2,
            "name": "Toyota Avalon",
            "stock": 2,
            "daily_rent": 500000
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve cars"
}
```

### GET /api/v1/cars/:id
Retrieve a specific car by ID.

**URL Parameters:**
- `id` (integer) - Car ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "name": "Toyota Camry",
        "stock": 2,
        "daily_rent": 500000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid car ID"
}

// 404 Not Found
{
    "error": "Car not found"
}
```

### POST /api/v1/cars
Create a new car.

**Request Body:**
```json
{
    "name": "BMW M3",
    "stock": 2,
    "daily_rent": 900000
}
```

**Field Requirements:**
- `name` (string, required) - Car model/name
- `stock` (integer, required) - Number of available cars (minimum 0)
- `daily_rent` (float, required) - Daily rental price (minimum 0)

**Success Response (201 Created):**
```json
{
    "data": {
        "no": 3,
        "name": "BMW M3",
        "stock": 2,
        "daily_rent": 900000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Missing required fields
{
    "error": "Key: 'Car.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}

// 400 Bad Request - Invalid stock (negative)
{
    "error": "Key: 'Car.Stock' Error:Field validation for 'Stock' failed on the 'min' tag"
}

// 500 Internal Server Error
{
    "error": "Failed to create car"
}
```

### PUT /api/v1/cars/:id
Update a specific car.

**URL Parameters:**
- `id` (integer) - Car ID

**Request Body:** (all fields optional)
```json
{
    "name": "BMW M3 - Updated",
    "stock": 2,
    "daily_rent": 900000
}
```

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 3,
        "name": "BMW M3 - Updated",
        "stock": 2,
        "daily_rent": 900000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid car ID
{
    "error": "Invalid car ID"
}

// 404 Not Found
{
    "error": "Car not found"
}

// 500 Internal Server Error
{
    "error": "Failed to update car"
}
```

### DELETE /api/v1/cars/:id
Delete a specific car.

**URL Parameters:**
- `id` (integer) - Car ID

**Success Response (200 OK):**
```json
{
    "message": "Car deleted successfully"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid car ID
{
    "error": "Invalid car ID"
}

// 404 Not Found
{
    "error": "Car not found"
}

// 500 Internal Server Error - Foreign key constraint
{
    "error": "Failed to delete car"
}
```

---

## Booking Endpoints

### GET /api/v1/bookings
Retrieve all bookings with customer and car details.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "customer_id": 3,
            "cars_id": 2,
            "start_rent": "2021-01-01T00:00:00Z",
            "end_rent": "2021-01-02T00:00:00Z",
            "total_cost": 1000000,
            "finished": true,
            "customer": {
                "no": 3,
                "name": "Hugo Fleming",
                "nik": "3372093912800",
                "phone_number": "081237123684"
            },
            "car": {
                "no": 2,
                "name": "Toyota Avalon",
                "stock": 2,
                "daily_rent": 500000
            }
        },
        {
            "no": 2,
            "customer_id": 1,
            "cars_id": 2,
            "start_rent": "2021-01-10T00:00:00Z",
            "end_rent": "2021-01-11T00:00:00Z",
            "total_cost": 1000000,
            "finished": true,
            "customer": {
                "no": 1,
                "name": "Wawan Hermawan",
                "nik": "3372093912739",
                "phone_number": "081237123682"
            },
            "car": {
                "no": 2,
                "name": "Toyota Avalon",
                "stock": 2,
                "daily_rent": 500000
            }
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve bookings"
}
```

### GET /api/v1/bookings/:id
Retrieve a specific booking by ID.

**URL Parameters:**
- `id` (integer) - Booking ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "customer_id": 3,
        "cars_id": 2,
        "start_rent": "2021-01-01T00:00:00Z",
        "end_rent": "2021-01-02T00:00:00Z",
        "total_cost": 1000000,
        "finished": true,
        "customer": {
            "no": 3,
            "name": "Hugo Fleming",
            "nik": "3372093912800",
            "phone_number": "081237123684"
        },
        "car": {
            "no": 2,
            "name": "Toyota Avalon",
            "stock": 2,
            "daily_rent": 500000
        }
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid booking ID"
}

// 404 Not Found
{
    "error": "Booking not found"
}
```

### POST /api/v1/bookings
Create a new booking.

**Request Body:**
```json
{
    "customer_id": 1,
    "cars_id": 1,
    "start_rent": "2025-07-05T00:00:00Z",
    "end_rent": "2025-07-07T00:00:00Z"
}
```

**Field Requirements:**
- `customer_id` (integer, required) - Must reference existing customer
- `cars_id` (integer, required) - Must reference existing car with stock > 0
- `start_rent` (datetime, required) - Must be before end_rent
- `end_rent` (datetime, required) - Must be after start_rent

**Success Response (201 Created):**
```json
{
    "data": {
        "no": 3,
        "customer_id": 1,
        "cars_id": 1,
        "start_rent": "2025-07-05T00:00:00Z",
        "end_rent": "2025-07-07T00:00:00Z",
        "total_cost": 1500000,
        "finished": false,
        "customer": {
            "no": 1,
            "name": "Wawan Hermawan",
            "nik": "3372093912739",
            "phone_number": "081237123682"
        },
        "car": {
            "no": 1,
            "name": "Toyota Camry",
            "stock": 1,
            "daily_rent": 500000
        }
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Missing required fields
{
    "error": "Key: 'Booking.CustomerID' Error:Field validation for 'CustomerID' failed on the 'required' tag"
}

// 400 Bad Request - Customer not found
{
    "error": "Customer not found"
}

// 400 Bad Request - Car not found
{
    "error": "Car not found"
}

// 400 Bad Request - Car not available
{
    "error": "Car is not available for booking"
}

// 400 Bad Request - Invalid date order
{
    "error": "Start date must be before end date"
}

// 400 Bad Request - Past date
{
    "error": "Start date cannot be in the past"
}

// 500 Internal Server Error
{
    "error": "Failed to create booking"
}
```

**Automatic Actions:**
- `total_cost` calculated as: (rental days) × (car daily rent)
- Car stock decremented by 1
- Booking marked as `finished: false`

### PUT /api/v1/bookings/:id
Update an existing booking.

**URL Parameters:**
- `id` (integer) - Booking ID

**Request Body:** (all fields optional)
```json
{
  "start_rent": "2025-07-06T00:00:00Z",
  "end_rent": "2025-07-09T00:00:00Z"
}
```

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 10,
        "customer_id": 1,
        "cars_id": 1,
        "start_rent": "2025-07-06T00:00:00Z",
        "end_rent": "2025-07-09T00:00:00Z",
        "total_cost": 2000000,
        "finished": false,
        "customer": {
            "no": 1,
            "name": "Wawan Hermawan",
            "nik": "3372093912739",
            "phone_number": "081237123682"
        },
        "car": {
            "no": 1,
            "name": "Toyota Camry",
            "stock": 1,
            "daily_rent": 500000
        }
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid booking ID
{
    "error": "Invalid booking ID"
}

// 404 Not Found
{
    "error": "Booking not found"
}

// 400 Bad Request - Cannot update finished booking
{
    "error": "Cannot update a finished booking"
}

// 400 Bad Request - Invalid dates
{
    "error": "Start date must be before end date"
}

// 500 Internal Server Error
{
    "error": "Failed to update booking"
}
```

**Notes:** 
- Cannot update finished bookings
- Date changes trigger cost recalculation

### DELETE /api/v1/bookings/:id
Delete a booking.

**URL Parameters:**
- `id` (integer) - Booking ID

**Success Response (200 OK):**
```json
{
    "message": "Booking deleted successfully"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid booking ID
{
    "error": "Invalid booking ID"
}

// 404 Not Found
{
    "error": "Booking not found"
}

// 400 Bad Request - Cannot delete finished booking
{
    "error": "Cannot delete a finished booking"
}

// 500 Internal Server Error
{
    "error": "Failed to delete booking"
}
```

**Automatic Actions:**
- Car stock restored (incremented by 1)
- Only applies to non-finished bookings

### PUT /api/v1/bookings/:id/finish
Mark a booking as finished.

**URL Parameters:**
- `id` (integer) - Booking ID

**Request Body:** None required

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 10,
        "customer_id": 1,
        "cars_id": 1,
        "start_rent": "2025-07-06T00:00:00Z",
        "end_rent": "2025-07-09T00:00:00Z",
        "total_cost": 2000000,
        "finished": true,
        "customer": {
            "no": 1,
            "name": "Wawan Hermawan",
            "nik": "3372093912739",
            "phone_number": "081237123682"
        },
        "car": {
            "no": 1,
            "name": "Toyota Camry",
            "stock": 2,
            "daily_rent": 500000
        }
    },
    "message": "Booking finished successfully"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid booking ID
{
    "error": "Invalid booking ID"
}

// 404 Not Found
{
    "error": "Booking not found"
}

// 400 Bad Request - Already finished
{
    "error": "Booking is already finished"
}

// 500 Internal Server Error
{
    "error": "Failed to finish booking"
}
```

**Automatic Actions:**
- Booking marked as `finished: true`
- Car stock restored (incremented by 1)

---

## Data Models

### Customer Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique customer identifier |
| `name` | string | ✅ | Not null | Customer's full name |
| `nik` | string | ✅ | Exactly 16 chars, Unique, Not null | National identification number |
| `phone_number` | string | ✅ | Max 15 chars, Not null | Customer's contact phone number |

### Car Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique car identifier |
| `name` | string | ✅ | Not null | Car model/name |
| `stock` | integer | ✅ | Min 0, Not null | Number of available cars |
| `daily_rent` | float | ✅ | Min 0, Not null | Daily rental price |

### Booking Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique booking identifier |
| `customer_id` | integer | ✅ | Foreign Key to Customer, Not null | Reference to customer |
| `cars_id` | integer | ✅ | Foreign Key to Car, Not null | Reference to car |
| `start_rent` | datetime | ✅ | Not null, >= current date | Rental start date and time |
| `end_rent` | datetime | ✅ | Not null, > start_rent | Rental end date and time |
| `total_cost` | float | - | Auto-calculated, Not null | Total cost for rental period |
| `finished` | boolean | - | Default: false | Completion status |

---

## Error Handling

### HTTP Status Codes

- `200 OK` - Successful GET, PUT, DELETE operations
- `201 Created` - Successful POST operations
- `400 Bad Request` - Invalid input data, validation errors, business rule violations
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Database or server errors

### Error Response Format

All error responses follow this format:
```json
{
  "error": "Descriptive error message"
}
```

### Common Error Messages

**Customer Errors:**
- `"Invalid customer ID"` - ID parameter is not a valid integer
- `"Customer not found"` - Customer doesn't exist
- `"Failed to create customer"` - Database error (often duplicate NIK)

**Car Errors:**
- `"Invalid car ID"` - ID parameter is not a valid integer
- `"Car not found"` - Car doesn't exist
- `"Failed to create car"` - Database error

**Booking Errors:**
- `"Invalid booking ID"` - ID parameter is not a valid integer
- `"Booking not found"` - Booking doesn't exist
- `"Car is not available for booking"` - Car stock is 0
- `"Cannot update a finished booking"` - Attempt to modify completed booking
- `"Cannot delete a finished booking"` - Attempt to delete completed booking
- `"Start date must be before end date"` - Invalid date range
- `"Start date cannot be in the past"` - Start date before current date

---

## Business Rules

### Stock Management
1. **Automatic Updates**: Car stock is automatically managed by booking operations
   - Creating booking: `stock = stock - 1`
   - Deleting booking: `stock = stock + 1`
   - Finishing booking: `stock = stock + 1`

2. **Availability**: Cars are bookable only when `stock > 0`

3. **Manual Updates**: Direct stock updates via car endpoints may interfere with booking logic

### Cost Calculation
- **Formula**: `total_cost = (rental_days) × (car.daily_rent)`
- **Day Calculation**: Includes both start and end dates (minimum 1 day)
- **Auto-Update**: Cost recalculated when booking dates change

### Validation Rules

**Customer Validation:**
- NIK must be exactly 16 characters and unique
- Phone number cannot exceed 15 characters
- All fields are required

**Car Validation:**
- Stock must be 0 or greater
- Daily rent must be 0 or greater
- All fields are required

**Booking Validation:**
- Customer and car must exist
- Car must be available (stock > 0)
- Start date cannot be in the past
- Start date must be before end date
- Cannot modify/delete finished bookings

### Foreign Key Constraints
- Customers with active bookings cannot be deleted
- Cars with active bookings cannot be deleted
- Booking operations validate referenced customers and cars exist

### Booking States
- `finished: false` - Active booking, car currently rented
- `finished: true` - Completed booking, car returned and stock restored
