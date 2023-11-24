package main

import (
	"broker-service/internal/application/adapter/api/routes"
	"broker-service/internal/service/broker_service"
	"broker-service/internal/service/handle_submission_service"
	"broker-service/internal/service/json_service"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	log.Printf("starting broker service on port %s\n", webPort)

	jsonService := json_service.New()
	brokerService := broker_service.New(jsonService)
	handleSubmissionService := handle_submission_service.New(jsonService)

	api_routes := routes.New(brokerService, handleSubmissionService)

	// define http service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	// start the service
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
