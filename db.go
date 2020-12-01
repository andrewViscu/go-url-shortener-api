package main

import (
	"context"
	"log"
	"os"
	"time"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

func DBConnect() *mongo.Client {
	ApplyURI := os.Getenv("DB_URI")
	opts := options.Client().ApplyURI(ApplyURI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v.\n Maybe the environment variables aren't set? Check them.", err)
	}
	log.Println("Connected to MongoDB!")
	return client
}