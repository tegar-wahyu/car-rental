package models

import "time"

type DriverIncentive struct {
	No        int        `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	BookingID int        `json:"booking_id" binding:"required" gorm:"column:booking_id;not null"`
	Incentive float64    `json:"incentive" binding:"required" gorm:"column:incentive;not null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	Booking Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID;references:No"`
}

func (DriverIncentive) TableName() string {
	return "driver_incentives"
}
