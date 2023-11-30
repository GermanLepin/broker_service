package connection

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURL = "mongodb://mongo:27017"
)

func ConnectToMongo() *mongo.Client {
	var counts int64

	for {
		clientOptions := options.Client().ApplyURI(mongoURL)
		clientOptions.SetAuth(options.Credential{
			Username: "mongousradmin",
			Password: "mongopassadmin",
		})

		connection, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Println("error connecting:", err)
			counts++
		} else {
			log.Println("connected to Mongo!")
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
