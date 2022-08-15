package services

import (
	"ArvanVoucher/models"
	"ArvanVoucher/requests"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type VoucherService struct {
	ServiceDB *gorm.DB
	Redis     *redis.Client
}

//AddVoucher creates new voucher
func (service VoucherService) AddVoucher(req requests.AddVoucherRequest) (*int64, error) {
	//check duplication from redis
	res, _ := service.Redis.Get(req.VoucherCode).Result()
	if len(res) > 0 {
		return nil, errors.New("VoucherAlreadyExists")
	}

	//check duplication from db
	var vch models.Vouchers
	service.ServiceDB.Where("voucher=? AND deleted_at IS NULL", req.VoucherCode).
		Find(&vch)
	if vch.Id > 0 {
		return nil, errors.New("VoucherAlreadyExists")
	}

	//add new voucher model
	trx := service.ServiceDB.Begin()
	expTime := req.ExpirationDate.Time
	voucher := models.Vouchers{
		Voucher:        req.VoucherCode,
		Amount:         req.Amount,
		Limit:          req.Limit,
		ExpirationDate: &expTime,
		CreatedAt:      time.Now(),
	}
	if e := trx.Create(&voucher).Error; e != nil {
		trx.Rollback()
		return nil, e
	}

	//add voucher with its data ro redis
	var duration time.Duration
	if req.ExpirationDate != nil {
		t := *req.ExpirationDate
		duration = t.Sub(time.Now())
	} else {
		duration = time.Hour * 24
	}

	queueData := requests.QueueRequest{
		Limit:  req.Limit,
		Amount: req.Amount,
		Exp:    duration,
	}
	redisData, _ := json.Marshal(queueData)
	if e := service.Redis.Set(req.VoucherCode, redisData, duration).Err(); e != nil {
		trx.Rollback()
		return nil, e
	}

	//add _count key for voucher to redis to guarantee atomic count for registrations
	if e := service.Redis.Set(req.VoucherCode+"_count", req.Limit, duration).Err(); e != nil {
		trx.Rollback()
		return nil, e
	}

	trx.Commit()
	return &voucher.Id, nil
}

//GetVoucherReport gives report for a specific voucher registrations
func (service VoucherService) GetVoucherReport(voucherCode string) ([]int64, error) {
	//iterate through redis for this voucher registrations
	var res []int64
	iter := REDIS.Scan(0, voucherCode+"#*", 0).Iterator()
	for iter.Next() {
		s := strings.Split(iter.Val(), "#")
		userId, e := strconv.ParseInt(s[1], 10, 64)
		if e != nil {
			return nil, e
		}
		res = append(res, userId)
	}

	//if redis fails to work, tries db
	if res == nil {
		var voucher models.Vouchers
		if e := DB.Where("voucher = ?", voucherCode).Find(&voucher).Error; e != nil {
			return nil, e
		}

		var voucherUser []models.VoucherUsers
		if e := DB.Where("voucher_id = ?", voucher.Id).Find(&voucherUser).Error; e != nil {
			return nil, e
		}

		for _, v := range voucherUser {
			res = append(res, v.UserId)
		}
	}

	return res, nil
}
