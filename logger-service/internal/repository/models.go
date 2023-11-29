package repository

import (
	"context"
	"log"
	"logger-service/internal/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *repository) InsertLogEntry(entry dto.LogEntry) error {
	collection := r.client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), dto.LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (r *repository) GetAll() ([]*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*dto.LogEntry

	for cursor.Next(ctx) {
		var item dto.LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (r *repository) GetOne(id string) (*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry dto.LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *repository) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateLogEntry(logEntry *dto.LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(logEntry.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", logEntry.Name},
				{"data", logEntry.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type repository struct {
	client *mongo.Client
}

func New(mongo *mongo.Client) *repository {
	return &repository{client: mongo}
}
