package models

type Membership struct {
	No             int     `json:"no" gorm:"primaryKey;column:no;autoIncrement"`
	MembershipName string  `json:"membership_name" binding:"required" gorm:"column:membership_name;not null"`
	Discount       float64 `json:"discount" binding:"required" gorm:"column:discount;not null"`
}

func (Membership) TableName() string {
	return "memberships"
}
