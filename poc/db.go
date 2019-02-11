package main

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/mongo/options"
)


const host string = "localhost"

// Cleanup will remove all mock data from the database.
func Cleanup(col string) {
	log.Println("Cleaning up MongoDB...")
	ctx, collection := GetMongo(col)

	_, err := collection.DeleteMany(ctx,
		bson.D{})
	if err != nil {
		log.Fatal(err)
	}
}

// GetMongo returns the session and a reference to the post collection.
func GetMongo(col string) (context.Context, *mongo.Collection) {
//  auth := options.Credential{
//    "MONGODB-CR",
//    nil,
//    "dt",
//    "admin",
//    "password",
//  }
	var options options.ClientOptions
	maxWait := time.Duration(5 * time.Second)
	options.SetConnectTimeout(maxWait)
//	options.SetAuth(auth)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, "mongodb://"+host+":27017", &options)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("minimalGraphql").Collection(col)
	return ctx, collection
}
