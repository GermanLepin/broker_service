package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/db/mongo/connection"
	"logger-service/internal/application/adapter/api/routes"
	"logger-service/internal/repository"
	"logger-service/internal/service/json_service"
	"logger-service/internal/service/logger_service"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

var client *mongo.Client

func main() {
	mongoClient, err := connection.ConnectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	logRepository := repository.New(mongoClient)

	jsonService := json_service.New()
	loggerService := logger_service.New(jsonService, logRepository)

	api_routes := routes.New(loggerService)

	log.Printf("starting logger service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
