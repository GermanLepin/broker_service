package main

import (
	"broker-service/db/rabbitmq/connection"
	"broker-service/internal/application/adapter/api/routes"
	"broker-service/internal/service/broker_service"
	"broker-service/internal/service/handle_submission_service"
	"broker-service/internal/service/json_service"
	"os"

	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	rabbitMQConn, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMQConn.Close()

	jsonService := json_service.New()
	brokerService := broker_service.New(jsonService)
	handleSubmissionService := handle_submission_service.New(rabbitMQConn, jsonService)

	api_routes := routes.New(brokerService, handleSubmissionService)

	log.Printf("starting broker service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
