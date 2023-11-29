package main

import (
	"authentication-service/db/postgres/connection"
	"authentication-service/internal/application/adapter/api/routes"
	"authentication-service/internal/repository"
	"authentication-service/internal/service/authentication_service"
	"authentication-service/internal/service/json_service"

	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"

func main() {
	conn := connection.StartDB()
	userRepository := repository.New(conn)

	jsonService := json_service.New()
	authenticationService := authentication_service.New(jsonService, userRepository)

	api_routes := routes.New(authenticationService)

	log.Printf("starting authentication service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
