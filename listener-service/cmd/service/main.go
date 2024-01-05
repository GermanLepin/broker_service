package main

import (
	"listener-service/db/rabbitmq/connection"
	"listener-service/internal/event"
	"log"
	"os"
)

func main() {
	rabbitMQConnection, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMQConnection.Close()

	log.Println("listing for and consuming RabbitMQ messages...")

	consumer, err := event.NewConnection(rabbitMQConnection)
	if err != nil {
		panic(err)
	}

	topics := []string{"log.INFO", "log.WARNING", "log.ERROR"}

	err = consumer.Listen(topics)
	if err != nil {
		log.Println(err)
	}
}
