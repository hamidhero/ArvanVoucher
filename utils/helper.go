package utils

import (
	"net/http"
	"time"
)

type Output struct {
	Timestamp  time.Time   `json:"TimeStamp"`
	Status     int         `json:"status"`
	SessionKey string      `json:"SessionKey"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      []Error     `json:"error"`
}

type Error struct {
	ErrorCode int
	ErrorMsg  string
}

func NewOutput() Output {
	output := Output{}
	output.Timestamp = time.Now()
	output.Status = http.StatusOK
	output.Message = "عملیات با موفقیت انجام شد"
	return output
}
