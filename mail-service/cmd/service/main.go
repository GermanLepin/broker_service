package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"mail-service/internal/application/adapter/api/routes"
	"mail-service/internal/dto"
	"mail-service/internal/service/json_service"
	"mail-service/internal/service/mail_service"

	"net/http"
)

const webPort = "80"

func main() {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	mail := dto.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	jsonService := json_service.New()
	mailService := mail_service.New(jsonService, mail)

	api_routes := routes.New(mailService)

	log.Println("starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
