package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/subosito/gotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var REDIS *redis.Client
var RABBIT *amqp.Connection
var CTX context.Context

func ConnectDB() error {
	if err := gotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}

	conn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"),
		os.Getenv("PASSWORD"), os.Getenv("DBNAME"))

	var e error
	DB, e = gorm.Open(postgres.Open(conn), &gorm.Config{})

	db, _ := DB.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour)

	return e
}

func ConnectRedis() error {
	REDIS = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, e := REDIS.Ping().Result()

	return e
}

func ConnectRabbitMq() error {
	var e error
	RABBIT, e = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if e != nil {
		return e
	}
	return nil
}
