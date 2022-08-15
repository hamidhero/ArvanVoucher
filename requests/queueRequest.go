package requests

import "time"

type AddToQueueRequest struct {
	Name   string      `json:"name"`
	Amount *int64      `json:"amount"`
	Data   interface{} `json:"data"`
}

type QueueRequest struct {
	Limit  int           `json:"limit"`
	Amount int64         `json:"amount"`
	Exp    time.Duration `json:"exp"`
}
