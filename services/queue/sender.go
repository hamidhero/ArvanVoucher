package queue

import (
	"ArvanVoucher/requests"
	"ArvanVoucher/services"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type QueueService struct {
	ServiceDB *gorm.DB
	Redis     *redis.Client
	Ctx       context.Context
}

func (service QueueService) AddToQueue(data requests.AddToQueueRequest) error {
	//decrease limit count and check if it is not exceeded. this part is atomic
	redisCount := service.Redis.Decr(data.Name + "_count").Val()
	if int(redisCount) < 1 {
		return errors.New("LimitExceeded")
	}

	//get voucher data
	res, e := service.Redis.Get(data.Name).Result()
	if e != nil {
		service.Redis.Incr(data.Name + "_count")
		return e
	}

	var queueData requests.QueueRequest
	e = json.Unmarshal([]byte(res), &queueData)
	if e != nil {
		service.Redis.Incr(data.Name + "_count")
		return e
	}

	//check user duplication for this voucher
	mobile, _ := json.Marshal(data.Data)
	key := data.Name + "#" + string(mobile)
	res, _ = service.Redis.Get(key).Result()
	if len(res) > 0 {
		service.Redis.Incr(data.Name + "_count")
		return errors.New("AlreadyRegistered")
	}

	//generate message to send on queue
	data.Amount = &queueData.Amount
	body, e := json.Marshal(data)
	if e != nil {
		service.Redis.Incr(data.Name + "_count")
		return e
	}

	//create rabbitMQ channel
	ch, e := services.RABBIT.Channel()
	if e != nil {
		service.Redis.Incr(data.Name + "_count")
		return e
	}
	defer ch.Close()

	//create an exchange in topic type
	e = ch.ExchangeDeclare(
		"vouchers_topic", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)

	//load balance requests
	routeKey := ".odd"
	if redisCount%2 == 0 {
		routeKey = ".even"
	}

	//publish message on queue
	e = ch.PublishWithContext(
		service.Ctx,
		"vouchers_topic",   // exchange
		data.Name+routeKey, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if e != nil {
		service.Redis.Incr(data.Name + "_count")
		return e
	}

	//add user to redis for this voucher to avoid duplication
	if e := service.Redis.Set(key, nil, queueData.Exp).Err(); e != nil {
		return e
	}

	//call a service to add to VoucherUsers model
	service.AddToVoucherUsersQueue(data)

	return nil
}

func (service QueueService) AddToVoucherUsersQueue(data requests.AddToQueueRequest) error {
	ch, e := services.RABBIT.Channel()
	if e != nil {
		return e
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
		return e
	}

	body, e := json.Marshal(data)
	if e != nil {
		return e
	}
	e = ch.PublishWithContext(
		service.Ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	return e
}
