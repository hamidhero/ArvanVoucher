package models

import "time"

type Vouchers struct {
	Id             int64      `gorm:"column:id;primary_key"`
	Voucher        string     `gorm:"column:voucher"`
	Amount         int64      `gorm:"column:amount"`
	Limit          int        `gorm:"column:limit"`
	ExpirationDate *time.Time `gorm:"column:expiration_date"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (Vouchers) TableName() string {
	return "Vouchers"
}
