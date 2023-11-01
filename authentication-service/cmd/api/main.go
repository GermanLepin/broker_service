package main

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

func main() {
	conn := connectToDB()
	if conn != nil {
		log.Panic("Cannot connect to Postgres")
	}

	config := NewConfig(conn)

	log.Printf("Starting authentication service on port %s\n", webPort)

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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func NewConfig(conn *sql.DB) *Config {
	return &Config{
		DB:     conn,
		Models: data.New(conn),
	}
}
