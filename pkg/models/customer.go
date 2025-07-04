package models

type Customer struct {
	No          int    `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	Name        string `json:"name" binding:"required" gorm:"column:name;not null"`
	NIK         string `json:"nik" binding:"required" gorm:"column:nik;not null;unique;size:16"`
	PhoneNumber string `json:"phone_number" binding:"required" gorm:"column:phone_number;not null;size:15"`
}

func (Customer) TableName() string {
	return "customers"
}
