package models

type BookingType struct {
	No          int    `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	BookingType string `json:"booking_type" binding:"required" gorm:"column:booking_type;not null"`
	Description string `json:"description" binding:"required" gorm:"column:description;not null"`
}

func (BookingType) TableName() string {
	return "booking_types"
}
