package connection

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost"
)

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64

	for {
		connection, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Println("rabbitMQ is not ready yet", err)
			counts++
		} else {
			log.Println("connected to rabbitMQ")
			return connection, nil
		}

		if counts > 10 {
			fmt.Println(err)
			return nil, err
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
