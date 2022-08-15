package main

import (
	"ArvanVoucher/services"
	"ArvanVoucher/services/queue"
	"log"
	"os"
)

func main() {
	var e error
	e = services.ConnectDB()
	if e != nil {
		log.Print("Error in db connection")
		println("Error in db connection")
	}

	e = services.ConnectRedis()
	if e != nil {
		log.Print("Error in redis connection")
		println("Error in redis connection")
	}

	e = services.ConnectRabbitMq()
	if e != nil {
		log.Print("Error in queue connection")
		println("Error in queue connection")
	}

	go func() {
		s := queue.QueueService{
			ServiceDB: services.DB,
		}
		s.ReceiveVoucherUsers()
	}()

	//time.Sleep(time.Second * 5)
	//
	//s := queue.QueueService{
	//	Rabbit:    services.RABBIT,
	//	ServiceDB: services.DB,
	//}
	//for i := 0; i < 10000000; i++ {
	//	sendData := requests.AddToQueueRequest{
	//		Name:   "t3",
	//		Data:   rand.Intn(9999999999),
	//	}
	//	s.AddToQueue(sendData)
	//}

	port := os.Args[1]
	GetRouter(port)
}
