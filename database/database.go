package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Uri string
}

// Connects to the MongoDB database that sits at the
// URI passed
//
//	d := common.Database{}
//	d.Connect("mongodb+srv://localhost:27017")
//
// The client is then returned on the database struct
//
// 	d.Client.Database("victoriam").Collection("users")
//
// Nothing is returned directly
func (d *Database) GetClient() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	options := options.Client().ApplyURI(d.Uri)
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging to MongoDB", err)
	}

	fmt.Println("Connected to MongoDB")
	return client, ctx, cancel
}

func (d *Database) GetCollection(collection string) (*mongo.Collection, context.Context, context.CancelFunc) {
	client, ctx, cancel := d.GetClient()
	return client.Database("victoriam").Collection(collection), ctx, cancel
}
