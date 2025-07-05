package models

import "time"

type Customer struct {
	No          int        `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	Name        string     `json:"name" binding:"required" gorm:"column:name;not null"`
	NIK         string     `json:"nik" binding:"required,len=16" gorm:"column:nik;not null;unique;size:16"`
	PhoneNumber string     `json:"phone_number" binding:"required,max=15" gorm:"column:phone_number;not null;size:15"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	MembershipID *int        `json:"membership_id" gorm:"column:membership_id"`
	Membership   *Membership `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:No"`
}

func (Customer) TableName() string {
	return "customers"
}
