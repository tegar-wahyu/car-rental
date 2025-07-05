package models

import (
	"time"
)

type Booking struct {
	No              int       `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	CustomerID      int       `json:"customer_id" binding:"required" gorm:"column:customer_id;not null"`
	CarsID          int       `json:"cars_id" binding:"required" gorm:"column:cars_id;not null"`
	StartRent       time.Time `json:"start_rent" binding:"required" gorm:"column:start_rent;not null"`
	EndRent         time.Time `json:"end_rent" binding:"required" gorm:"column:end_rent;not null"`
	TotalCost       float64   `json:"total_cost" gorm:"column:total_cost;not null"`
	Finished        bool      `json:"finished" gorm:"column:finished;default:false"`
	Discount        float64   `json:"discount" gorm:"column:discount;default:0"`
	BookingTypeID   int       `json:"booking_type_id" binding:"required" gorm:"column:booking_type_id;not null"`
	DriverID        *int      `json:"driver_id" gorm:"column:driver_id"`
	TotalDriverCost float64   `json:"total_driver_cost" gorm:"column:total_driver_cost;default:0"`

	Customer    Customer    `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:No" binding:"-"`
	Car         Car         `json:"car,omitempty" gorm:"foreignKey:CarsID;references:No" binding:"-"`
	Driver      *Driver     `json:"driver,omitempty" gorm:"foreignKey:DriverID;references:No" binding:"-"`
	BookingType BookingType `json:"booking_type,omitempty" gorm:"foreignKey:BookingTypeID;references:No" binding:"-"`
}

type BookingUpdate struct {
	StartRent       *time.Time `json:"start_rent,omitempty"`
	EndRent         *time.Time `json:"end_rent,omitempty"`
	TotalCost       *float64   `json:"total_cost,omitempty"`
	Finished        *bool      `json:"finished,omitempty"`
	Discount        *float64   `json:"discount,omitempty"`
	TotalDriverCost *float64   `json:"total_driver_cost,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}
