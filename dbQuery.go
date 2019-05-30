package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func countItems(key string, value string) (int64, error) {
	ctx, collection := GetMongo("person")
	cnt, err := collection.CountDocuments(
		ctx,
		bson.D{
			{key, value},
		},
	)
	if err != nil {
		log.Println(err)
	}
	return cnt, err
}
