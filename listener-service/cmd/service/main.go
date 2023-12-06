package main

import (
	"listener-service/db/rabbitmq/connection"
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

}
