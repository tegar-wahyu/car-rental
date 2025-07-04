# Car Rental API Documentation

This document provides comprehensive documentation for all API endpoints in the car rental management system.

## Table of Contents

- [Overview](#overview)
- [Base URL](#base-url)
- [Health Check](#health-check)
- [Customer Endpoints](#customer-endpoints)
- [Car Endpoints](#car-endpoints)
- [Booking Endpoints](#booking-endpoints)
- [Membership Endpoints](#membership-endpoints)
- [Driver Endpoints](#driver-endpoints)
- [Booking Type Endpoints](#booking-type-endpoints)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Business Rules](#business-rules)

## Overview

The Car Rental API provides complete CRUD operations for managing customers, cars, bookings, memberships, and drivers. The API features automatic stock management, cost calculation with membership discounts, driver assignments with cost calculation, and comprehensive validation with constraint-based error handling.

### Key Features
- **Automatic Stock Management** - Car inventory automatically updated on booking operations
- **Cost Calculation** - Total costs computed based on rental duration, daily rates, membership discounts, and driver costs
- **Membership System** - Customer membership integration with discount calculation
- **Driver Service** - Car & Driver rental option with driver cost calculation
- **Booking Types** - Support for different booking types (Car Only, Car & Driver)
- **Constraint-Based Validation** - Advanced referential integrity checking with detailed error responses
- **Data Validation** - Input validation with detailed error messages
- **Relationship Management** - Comprehensive foreign key handling and constraints

## Base URL

```
http://localhost:8080
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

### GET /customers
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

### GET /customers/:id
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

### POST /customers
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

### PUT /customers/:id
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

### DELETE /customers/:id
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

// 400 Bad Request - Customer has active bookings
{
    "error": "Cannot delete customer with active bookings. Please finish or cancel active bookings first.",
    "entity_type": "customer",
    "entity_id": 1,
    "constraint": "active_bookings",
    "details": {
        "active_bookings": 2,
        "total_bookings": 5
    }
}

// 400 Bad Request - Customer has booking history  
{
    "error": "Cannot delete customer with booking history. Customer has completed bookings in the system.",
    "entity_type": "customer", 
    "entity_id": 1,
    "constraint": "booking_history",
    "details": {
        "active_bookings": 0,
        "total_bookings": 3
    }
}
```

### PUT /customers/:id/subscribe/:membership_id
Subscribe a customer to a membership plan.

**URL Parameters:**
- `id` (integer) - Customer ID
- `membership_id` (integer) - Membership plan ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "name": "Wawan Hermawan",
        "nik": "3372093912739",
        "phone_number": "081237123682",
        "membership_id": 2,
        "membership": {
            "no": 2,
            "membership_name": "Silver",
            "discount": 7
        }
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid customer ID
{
    "error": "Invalid customer ID"
}

// 404 Not Found - Customer doesn't exist
{
    "error": "Customer not found"
}

// 400 Bad Request - Invalid membership ID
{
    "error": "Invalid membership ID"
}

// 400 Bad Request - Membership doesn't exist
{
    "error": "Membership not found"
}

// 500 Internal Server Error
{
    "error": "Failed to subscribe customer to membership"
}
```

### DELETE /customers/:id/unsubscribe
Remove a customer's membership subscription.

**URL Parameters:**
- `id` (integer) - Customer ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 2,
        "name": "Philip Walker",
        "nik": "3372093912785",
        "phone_number": "081237123683",
        "membership_id": null
    },
    "message": "Successfully unsubscribed from membership"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid customer ID
{
    "error": "Invalid customer ID"
}

// 404 Not Found - Customer doesn't exist
{
    "error": "Customer not found"
}

// 500 Internal Server Error
{
    "error": "Failed to unsubscribe customer from membership"
}
```

---

## Car Endpoints

### GET /cars
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

### GET /cars/:id
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

### POST /cars
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

### PUT /cars/:id
Update an existing car.

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

### DELETE /cars/:id
Delete a car.

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

// 400 Bad Request - Car has active bookings
{
    "error": "Cannot delete car with active bookings. Please finish or cancel active bookings first.",
    "entity_type": "car",
    "entity_id": 1,
    "constraint": "active_bookings",
    "details": {
        "active_bookings": 1,
        "total_bookings": 5
    }
}

// 400 Bad Request - Car has booking history
{
    "error": "Cannot delete car with booking history. Car has completed bookings in the system.",
    "entity_type": "car",
    "entity_id": 1,
    "constraint": "booking_history", 
    "details": {
        "active_bookings": 0,
        "total_bookings": 3
    }
}

// 500 Internal Server Error
{
    "error": "Failed to delete car from database"
}
```

---

## Booking Endpoints

### GET /bookings
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
            "discount": 0,
            "booking_type_id": 1,
            "driver_id": null,
            "total_driver_cost": 0,
            "customer": {
                "no": 3,
                "name": "Hugo Fleming",
                "nik": "3372093912800",
                "phone_number": "081237123684",
                "membership_id": null
            },
            "car": {
                "no": 2,
                "name": "Toyota Avalon",
                "stock": 2,
                "daily_rent": 500000
            },
            "booking_type": {
                "no": 1,
                "booking_type": "Car Only",
                "description": "Rent Car only"
            }
        },
        {
            "no": 2,
            "customer_id": 11,
            "cars_id": 2,
            "start_rent": "2021-01-10T00:00:00Z",
            "end_rent": "2021-01-11T00:00:00Z",
            "total_cost": 1000000,
            "finished": true,
            "discount": 40000,
            "booking_type_id": 1,
            "driver_id": null,
            "total_driver_cost": 0,
            "customer": {
                "no": 11,
                "name": "Damien Kaufman",
                "nik": "3372093913202",
                "phone_number": "081237123692",
                "membership_id": 1,
                "membership": {
                    "no": 1,
                    "membership_name": "Bronze",
                    "discount": 4
                }
            },
            "car": {
                "no": 2,
                "name": "Toyota Avalon",
                "stock": 2,
                "daily_rent": 500000
            },
            "booking_type": {
                "no": 1,
                "booking_type": "Car Only",
                "description": "Rent Car only"
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

### GET /bookings/:id
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

### POST /bookings
Create a new booking.

**Request Body:**
```json
{
    "customer_id": 1,
    "cars_id": 1,
    "start_rent": "2025-07-05T00:00:00Z",
    "end_rent": "2025-07-07T00:00:00Z",
    "booking_type_id": 2,
    "driver_id": 2
}
```

**Field Requirements:**
- `customer_id` (integer, required) - Must reference existing customer
- `cars_id` (integer, required) - Must reference existing car with stock > 0
- `booking_type_id` (integer, required) - Must reference existing booking type
- `start_rent` (datetime, required) - Must be before end_rent
- `end_rent` (datetime, required) - Must be after start_rent
- `driver_id` (integer, optional) - Must reference existing driver (required for "Car & Driver" booking type)

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
        "discount": 0,
        "booking_type_id": 2,
        "driver_id": 2,
        "total_driver_cost": 405000,
        "customer": {
            "no": 1,
            "name": "Wawan Hermawan",
            "nik": "3372093912739",
            "phone_number": "081237123682",
            "membership_id": null
        },
        "car": {
            "no": 1,
            "name": "Toyota Camry",
            "stock": 1,
            "daily_rent": 500000
        },
        "driver": {
            "no": 2,
            "name": "Halsey Quinn",
            "nik": "3220132938293",
            "phone_number": "081992048713",
            "daily_cost": 135000
        },
        "booking_type": {
            "no": 2,
            "booking_type": "Car & Driver",
            "description": "Rent Car and a Driver"
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

// 400 Bad Request - Booking type not found
{
    "error": "Booking type not found"
}

// 400 Bad Request - Driver not found
{
    "error": "Driver not found"
}

// 400 Bad Request - Driver validation for Car & Driver booking
{
    "error": "Driver must be assigned for 'Car & Driver' booking type"
}

// 400 Bad Request - Driver validation for Car Only booking
{
    "error": "Driver can only be assigned for 'Car & Driver' booking type"
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
- `discount` calculated based on customer membership
- `total_driver_cost` calculated as: (rental days) × (driver daily cost) if driver assigned
- Car stock decremented by 1
- Booking marked as `finished: false`

### PUT /bookings/:id
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

### DELETE /bookings/:id
Delete a booking.

**URL Parameters:**
- `id` (integer) - Booking ID

**Success Response (200 OK):**
```json
{
    "details": {
        "deleted_booking_id": 11,
        "restored_car_stock": 3
    },
    "message": "Booking deleted successfully and car stock restored"
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
    "error": "Cannot delete finished booking. Finished bookings are kept for historical records.",
    "entity_type": "booking",
    "entity_id": 1,
    "constraint": "finished_booking",
    "details": {
        "booking_id": 1,
        "finished": true
    }
}

// 500 Internal Server Error
{
    "error": "Failed to delete booking from database"
}

// 500 Internal Server Error
{
    "error": "Failed to restore car stock"
}
```

**Automatic Actions:**
- Car stock restored (incremented by 1)
- Only applies to non-finished bookings

### PUT /bookings/:id/finish
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

## Membership Endpoints

**Note:** Membership endpoints are read-only. Memberships are managed by the system and cannot be created, updated, or deleted via the API. Customers can subscribe to existing memberships using the customer subscription endpoints.

### GET /memberships
Retrieve all memberships.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "membership_name": "Bronze",
            "discount": 4
        },
        {
            "no": 2,
            "membership_name": "Silver",
            "discount": 7
        },
        {
            "no": 3,
            "membership_name": "Gold",
            "discount": 15
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve memberships"
}
```

### GET /memberships/:id
Retrieve a specific membership by ID.

**URL Parameters:**
- `id` (integer) - Membership ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "membership_name": "Bronze",
        "discount": 4
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid membership ID"
}

// 404 Not Found - Membership doesn't exist
{
    "error": "Membership not found"
}
```

---

## Driver Endpoints

### GET /drivers
Retrieve all drivers.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "name": "Stanley Baxter",
            "nik": "3220132938273",
            "phone_number": "81992048712",
            "daily_cost": 150000
        },
        {
            "no": 2,
            "name": "Halsey Quinn",
            "nik": "3220132938293",
            "phone_number": "081992048713",
            "daily_cost": 135000
        },
        {
            "no": 3,
            "name": "Kingsley Alvarez",
            "nik": "3220132938313",
            "phone_number": "081992048714",
            "daily_cost": 150000
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve drivers"
}
```

### GET /drivers/:id
Retrieve a specific driver by ID.

**URL Parameters:**
- `id` (integer) - Driver ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "name": "Stanley Baxter",
        "nik": "3220132938273",
        "phone_number": "81992048712",
        "daily_cost": 150000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid driver ID"
}

// 404 Not Found - Driver doesn't exist
{
    "error": "Driver not found"
}
```

### POST /drivers
Create a new driver.

**Request Body:**
```json
{
    "name": "Ahmad Driver",
    "nik": "9876543210987654",
    "phone_number": "0821234567",
    "daily_cost": 150000.0
}
```

**Field Requirements:**
- `name` (string, required) - Driver's full name
- `nik` (string, required) - National identification number (exactly 16 characters, must be unique)
- `phone_number` (string, required) - Driver's contact phone number (max 15 characters)
- `daily_cost` (float, required) - Driver's daily cost (minimum 0)

**Success Response (201 Created):**
```json
{
    "data": {
        "no": 4,
        "name": "Ahmad Driver",
        "nik": "9876543210987654",
        "phone_number": "0821234567",
        "daily_cost": 150000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Missing required fields
{
    "error": "Key: 'Driver.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}

// 400 Bad Request - Invalid NIK length
{
    "error": "Key: 'Driver.NIK' Error:Field validation for 'NIK' failed on the 'required' tag"
}

// 400 Bad Request - Invalid daily cost
{
    "error": "Key: 'Driver.DailyCost' Error:Field validation for 'DailyCost' failed on the 'min' tag"
}

// 500 Internal Server Error - Duplicate NIK
{
    "error": "Failed to create driver"
}
```

### PUT /drivers/:id
Update an existing driver.

**URL Parameters:**
- `id` (integer) - Driver ID

**Request Body:**
```json
{
    "name": "Ahmad Driver Updated",
    "nik": "9876543210987654",
    "phone_number": "0821234568",
    "daily_cost": 175000.0
}
```

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 8,
        "name": "Ahmad Driver Updated",
        "nik": "9876543210987654",
        "phone_number": "0821234568",
        "daily_cost": 175000
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid driver ID
{
    "error": "Invalid driver ID"
}

// 404 Not Found
{
    "error": "Driver not found"
}

// 500 Internal Server Error
{
    "error": "Failed to update driver"
}
```

### DELETE /drivers/:id
Delete a driver.

**URL Parameters:**
- `id` (integer) - Driver ID

**Success Response (200 OK):**
```json
{
    "message": "Driver deleted successfully"
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid driver ID
{
    "error": "Invalid driver ID"
}

// 404 Not Found
{
    "error": "Driver not found"
}

// 400 Bad Request - Driver has active bookings
{
    "error": "Cannot delete driver with active bookings. Please finish or cancel active bookings first.",
    "entity_type": "driver",
    "entity_id": 1,
    "constraint": "active_bookings",
    "details": {
        "active_bookings": 2,
        "total_bookings": 10
    }
}

// 400 Bad Request - Driver has booking history
{
    "error": "Cannot delete driver with booking history. Driver has completed bookings in the system.",
    "entity_type": "driver",
    "entity_id": 1,
    "constraint": "booking_history",
    "details": {
        "active_bookings": 0,
        "total_bookings": 5
    }
}

// 500 Internal Server Error
{
    "error": "Failed to delete driver from database"
}
```

### GET /drivers/:id/incentives
Retrieve driver incentives for a specific driver.

**URL Parameters:**
- `id` (integer) - Driver ID

**Success Response (200 OK):**
```json

// Driver with incentives
{
    "data": [
        {
            "no": 1,
            "booking_id": 6,
            "incentive": 40000,
            "booking": {
                "no": 6,
                "customer_id": 12,
                "cars_id": 14,
                "start_rent": "2021-02-16T00:00:00Z",
                "end_rent": "2021-02-16T00:00:00Z",
                "total_cost": 800000,
                "finished": true,
                "discount": 56000,
                "booking_type_id": 2,
                "driver_id": 2,
                "total_driver_cost": 135000,
                "customer": {
                    "no": 12,
                    "name": "Ayesha Richardson",
                    "nik": "3372093913257",
                    "phone_number": "081237123693",
                    "membership_id": 2
                },
                "car": {
                    "no": 14,
                    "name": "Mitsubishi Pajero Sport",
                    "stock": 5,
                    "daily_rent": 800000
                },
                "booking_type": {
                    "no": 0,
                    "booking_type": "",
                    "description": ""
                }
            }
        }
    ],
    "total_incentives": 40000
}

// Driver without incentives
{
    "data": [],
    "total_incentives": 0
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid driver ID
{
    "error": "Invalid driver ID"
}

// 404 Not Found
{
    "error": "Driver not found"
}

// 500 Internal Server Error
{
    "error": "Failed to retrieve driver incentives"
}
```

---

## Booking Type Endpoints

**Note:** Booking type endpoints are read-only and accessed via `/bookings/types`. Booking types are managed by the system and cannot be created, updated, or deleted via the API.

### GET /bookings/types
Retrieve all booking types.

**Success Response (200 OK):**
```json
{
    "data": [
        {
            "no": 1,
            "booking_type": "Car Only",
            "description": "Rent Car only"
        },
        {
            "no": 2,
            "booking_type": "Car & Driver",
            "description": "Rent Car and a Driver"
        }
    ]
}
```

**Error Responses:**
```json
// 500 Internal Server Error
{
    "error": "Failed to retrieve booking types"
}
```

### GET /bookings/types/:id
Retrieve a specific booking type by ID.

**URL Parameters:**
- `id` (integer) - Booking Type ID

**Success Response (200 OK):**
```json
{
    "data": {
        "no": 1,
        "booking_type": "Car Only",
        "description": "Rent Car only"
    }
}
```

**Error Responses:**
```json
// 400 Bad Request - Invalid ID format
{
    "error": "Invalid booking type ID"
}

// 404 Not Found - Booking type doesn't exist
{
    "error": "Booking type not found"
}
```

---

## Data Models

### Customer Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique customer identifier |
| `name` | string | ✅ | Not null | Customer's full name |
| `nik` | string | ✅ | Exactly 16 chars, Unique, Not null | National identification number |
| `phone_number` | string | ✅ | Max 15 chars, Not null | Customer's contact phone number |
| `membership_id` | integer | - | Foreign Key to Membership, Nullable | Reference to membership plan |
| `membership` | object | - | Populated when preloaded | Membership details object |

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
| `booking_type_id` | integer | ✅ | Foreign Key to BookingType, Not null | Reference to booking type |
| `driver_id` | integer | - | Foreign Key to Driver, Nullable | Reference to driver (required for Car & Driver) |
| `start_rent` | datetime | ✅ | Not null, >= current date | Rental start date and time |
| `end_rent` | datetime | ✅ | Not null, > start_rent | Rental end date and time |
| `total_cost` | float | - | Auto-calculated, Not null | Total cost for rental period |
| `discount` | float | - | Auto-calculated, Default: 0 | Membership discount amount |
| `total_driver_cost` | float | - | Auto-calculated, Default: 0 | Total driver cost for rental period |
| `finished` | boolean | - | Default: false | Completion status |

### Membership Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique membership identifier |
| `membership_name` | string | ✅ | Not null, Unique | Membership name |
| `discount` | float | ✅ | Min 0, Max 100, Not null | Discount percentage |

### Driver Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique driver identifier |
| `name` | string | ✅ | Not null | Driver's full name |
| `nik` | string | ✅ | Exactly 16 chars, Unique, Not null | National identification number |
| `phone_number` | string | ✅ | Max 15 chars, Not null | Driver's contact phone number |
| `daily_cost` | float | ✅ | Min 0, Not null | Driver's daily cost |

### Booking Type Model

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `no` | integer | - | Auto-generated, Primary Key | Unique booking type identifier |
| `booking_type` | string | ✅ | Not null | Booking type name (e.g., "Car Only", "Car & Driver") |
| `description` | string | ✅ | Not null | Detailed description of the booking type |

---

## Error Handling

### HTTP Status Codes

- `200 OK` - Successful GET, PUT, DELETE operations
- `201 Created` - Successful POST operations
- `400 Bad Request` - Invalid input data, validation errors, business rule violations
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Database or server errors

### Error Response Format

Standard error responses follow this format:
```json
{
  "error": "Descriptive error message"
}
```

Constraint-based error responses (for deletion operations) follow this enhanced format:
```json
{
  "error": "Descriptive error message with suggested action",
  "entity_type": "customer|car|driver|booking",
  "entity_id": 123,
  "constraint": "active_bookings|booking_history|finished_booking",
  "details": {
    "active_bookings": 2,
    "total_bookings": 5,
    "additional_context": "value"
  }
}
```

### Common Error Messages

**Customer Errors:**
- `"Invalid customer ID"` - ID parameter is not a valid integer
- `"Customer not found"` - Customer doesn't exist
- `"Failed to create customer"` - Database error (often duplicate NIK)
- `"Cannot delete customer with active bookings..."` - Customer has active bookings
- `"Cannot delete customer with booking history..."` - Customer has completed bookings

**Car Errors:**
- `"Invalid car ID"` - ID parameter is not a valid integer
- `"Car not found"` - Car doesn't exist
- `"Failed to create car"` - Database error
- `"Cannot delete car with active bookings..."` - Car has active bookings
- `"Cannot delete car with booking history..."` - Car has completed bookings

**Booking Errors:**
- `"Invalid booking ID"` - ID parameter is not a valid integer
- `"Booking not found"` - Booking doesn't exist
- `"Car is not available for booking"` - Car stock is 0
- `"Cannot update a finished booking"` - Attempt to modify completed booking
- `"Cannot delete finished booking..."` - Attempt to delete completed booking
- `"Start date must be before end date"` - Invalid date range
- `"Start date cannot be in the past"` - Start date before current date
- `"Booking type not found"` - Invalid booking type reference
- `"Driver not found"` - Invalid driver reference
- `"Driver must be assigned for 'Car & Driver' booking type"` - Missing driver for Car & Driver
- `"Driver can only be assigned for 'Car & Driver' booking type"` - Invalid driver assignment

**Membership Errors:**
- `"Invalid membership ID"` - ID parameter is not a valid integer
- `"Membership not found"` - Membership doesn't exist

**Driver Errors:**
- `"Invalid driver ID"` - ID parameter is not a valid integer
- `"Driver not found"` - Driver doesn't exist
- `"Failed to create driver"` - Database error (often duplicate NIK)
- `"Cannot delete driver with active bookings..."` - Driver has active bookings
- `"Cannot delete driver with booking history..."` - Driver has completed bookings

**Booking Type Errors:**
- `"Invalid booking type ID"` - ID parameter is not a valid integer
- `"Booking type not found"` - Booking type doesn't exist

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
- **Base Formula**: `total_cost = (rental_days) × (car.daily_rent)`
- **Membership Discount**: `discount = total_cost × (customer.membership.discount / 100)` if customer has membership
- **Driver Cost**: `total_driver_cost = (rental_days) × (driver.daily_cost)` if driver assigned
- **Day Calculation**: Includes both start and end dates (minimum 1 day)
- **Auto-Update**: Cost, discount, and driver cost recalculated when booking dates change

### Booking Type Validation
- **Car Only**: Driver assignment not allowed (`driver_id` must be null)
- **Car & Driver**: Driver assignment required (`driver_id` must not be null)
- Driver must exist and be valid when assigned

### Validation Rules

**Customer Validation:**
- NIK must be exactly 16 characters and unique
- Phone number cannot exceed 15 characters
- All fields are required
- Membership subscription is optional

**Car Validation:**
- Stock must be 0 or greater
- Daily rent must be 0 or greater
- All fields are required

**Driver Validation:**
- NIK must be exactly 16 characters and unique
- Phone number cannot exceed 15 characters
- Daily cost must be 0 or greater
- All fields are required

**Booking Validation:**
- Customer, car, and booking type must exist
- Car must be available (stock > 0)
- Start date cannot be in the past
- Start date must be before end date
- Cannot modify/delete finished bookings
- Driver validation based on booking type

### Foreign Key Constraints & Referential Integrity
- **Customers**: Cannot be deleted if they have active bookings; deletion blocked if they have any booking history
- **Cars**: Cannot be deleted if they have active bookings; deletion blocked if they have any booking history  
- **Drivers**: Cannot be deleted if they have active bookings; deletion blocked if they have any booking history
- **Bookings**: Cannot be deleted if marked as finished (kept for historical records)
- **Memberships & Booking Types**: Read-only via API, managed by system
- All booking operations validate that referenced customers, cars, drivers, and booking types exist
- Constraint violations return detailed error responses with entity information and suggested actions

### Booking States
- `finished: false` - Active booking, car currently rented
- `finished: true` - Completed booking, car returned and stock restored
