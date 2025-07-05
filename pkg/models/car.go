package models

import "time"

type Car struct {
	No        int        `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	Name      string     `json:"name" binding:"required" gorm:"column:name;not null"`
	Stock     int        `json:"stock" binding:"required,min=0" gorm:"column:stock;not null"`
	DailyRent float64    `json:"daily_rent" binding:"required,min=0" gorm:"column:daily_rent;not null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Car) TableName() string {
	return "cars"
}
