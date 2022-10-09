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

// Creates an index
//
//	// Single Index
//	d.SetIndex("users", bson.D{{Key: "accessToken", Value: 1}})
//	// Compound Index
//	d.SetIndex("products", bson.D{{Key: "tags", value: 1}, {Key: "productType", Value: 1}})
//	// Test Index
//	d.SetIndex("products", bson.D{{Key: "description", Value: "text"}})
//
// Returns an error when the create index fails
func (d *Database) SetIndex(collection string, model interface{}) error {
	coll, ctx, cancel := d.GetCollection(collection)
	defer cancel()

	// Set the index if it doesn't currently exist
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: model})
	return err
}
