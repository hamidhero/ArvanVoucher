package queue

import (
	"ArvanVoucher/models"
	"ArvanVoucher/requests"
	"ArvanVoucher/services"
	"encoding/json"
	"strconv"
	"time"
)

//ReceiveVoucherUsers consumer for updating VoucherUsers model
func (service QueueService) ReceiveVoucherUsers() {
	ch, e := services.RABBIT.Channel()
	if e != nil {
		return
	}
	defer ch.Close()

	q, e := ch.QueueDeclare(
		"VoucherUsers", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if e != nil {
		return
	}

	msgs, e := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if e != nil {
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var data requests.AddToQueueRequest
			e := json.Unmarshal(d.Body, &data)
			if e == nil {
				//finding voucher record
				var voucher models.Vouchers
				service.ServiceDB.Where("voucher=?", data.Name).Find(&voucher)
				mobileByte, _ := json.Marshal(data.Data)
				mobile, _ := strconv.ParseInt(string(mobileByte), 10, 64)

				//creating new row in VoucherUsers
				voucherUser := models.VoucherUsers{
					VoucherId: voucher.Id,
					UserId:    mobile,
					CreatedAt: time.Now(),
				}
				e := service.ServiceDB.Create(&voucherUser).Error
				if e == nil {
					d.Ack(false)
				}
			}
		}
	}()

	<-forever
}
