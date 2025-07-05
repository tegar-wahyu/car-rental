package models

import "time"

type Driver struct {
	No          int        `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	Name        string     `json:"name" binding:"required" gorm:"column:name;not null"`
	NIK         string     `json:"nik" binding:"required" gorm:"column:nik;not null;unique;size:16"`
	PhoneNumber string     `json:"phone_number" binding:"required" gorm:"column:phone_number;not null;size:15"`
	DailyCost   float64    `json:"daily_cost" binding:"required,min=0" gorm:"column:daily_cost;not null"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Driver) TableName() string {
	return "drivers"
}
