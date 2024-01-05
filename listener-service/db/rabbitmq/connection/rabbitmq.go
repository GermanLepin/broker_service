package connection

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL = "amqp://guest:guest@rabbitmq"
)

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	var (
		counts  int64
		backOff = 1 * time.Second
	)

	for {
		connection, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Println("rabbitMQ is not ready yet", err)
			counts++
		} else {
			log.Println("connected to rabbitMQ")
			return connection, nil
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
}
