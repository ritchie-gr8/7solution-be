package databases

import (
	"context"
	"log"
	"time"

	"github.com/ritchie-gr8/7solution-be/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConnect(cfg config.IDBConfig) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Url()).SetMaxPoolSize(uint64(cfg.MaxPoolSize()))
	db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	log.Printf("Connected to MongoDB with max pool size: %d", cfg.MaxPoolSize())
	return db
}

func DbDisconnect(db *mongo.Client) {
	if db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Disconnect(ctx); err != nil {
		log.Fatalf("failed to disconnect from db: %v", err)
	}

	log.Println("Disconnected from MongoDB")
}
