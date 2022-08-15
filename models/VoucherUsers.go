package models

import "time"

type VoucherUsers struct {
	Id        int64      `gorm:"column:id;primary_key"`
	VoucherId int64      `gorm:"column:voucher_id"`
	UserId    int64      `gorm:"column:user_id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (VoucherUsers) TableName() string {
	return "VoucherUsers"
}
