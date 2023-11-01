package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	config := NewConfig()

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: config.routes(),
	}

	// start the service
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}
