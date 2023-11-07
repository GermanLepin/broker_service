package main

import (
	"authentication-service/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose"
)

const webPort = "80"

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("cannot connect to Postgres")
	}

	if err := goose.Up(conn, "/var"); err != nil {
		log.Panic("cannot do the migrations")
	}

	config := NewConfig(conn)

	log.Printf("starting authentication service on port %s\n", webPort)

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

	var counts int64

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres is not ready yet")
			counts++
		} else {
			log.Println("connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

type Config struct {
	DB     *sql.DB
	Models repository.Models
}

func NewConfig(conn *sql.DB) *Config {
	return &Config{
		DB:     conn,
		Models: repository.New(conn),
	}
}
