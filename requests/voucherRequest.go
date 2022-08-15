package requests

import (
	"encoding/json"
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

func (ct *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	s = strings.Replace(s, "\\", "", 1)
	if s == "null" {
		return
	}
	ct.Time, err = dateparse.ParseAny(s)
	return
}

type AddVoucherRequest struct {
	VoucherCode    string    `json:"voucher_code"`
	Amount         int64     `json:"amount"`
	Limit          int       `json:"limit"`
	ExpirationDate *DateTime `json:"expiration_date"`
}

type RegisterVoucherRequest struct {
	Mobile      json.Number `json:"mobile"`
	VoucherCode string      `json:"voucher_code"`
}
