package main

import (
	"context"
	//  "fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var host string
var dbName string

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}
	host = os.Getenv("MONGODB_URI")
	log.Println(`MONGO_URI: ` + host)
	// Get the string without protocol
	host_splited := strings.Split(host, "//")
	// Get the dbname, if there is any
	host_splited = strings.Split(host_splited[len(host_splited)-1], "/")
	if len(host_splited) == 2 {
		dbName = host_splited[len(host_splited)-1]
	} else {
		dbName = "digitaltmc"
	}
	log.Println(`dbName: ` + dbName)
}

// Cleanup will remove all mock data from the database.
func Cleanup(col string) {
	log.Println("Cleaning up MongoDB...")
	ctx, collection := GetMongo(col)

	_, err := collection.DeleteMany(ctx,
		bson.D{})
	if err != nil {
		log.Println(err)
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
	// 	var options options.ClientOptions
	// 	maxWait := time.Duration(5 * time.Second)
	// 	options.SetConnectTimeout(maxWait)
	// //	options.SetAuth(auth)
	//
	// 	ctx := context.Background()

	// client, err := mongo.Connect(ctx, host, &options)

	// Use one-stop solution instead.
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))

	if err != nil {
		log.Println(err)
	}
	collection := client.Database(dbName).Collection(col)
	return ctx, collection
}
