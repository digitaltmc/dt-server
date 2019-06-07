package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (_ *Resolver) Book(args struct {
	Token string
	Date  string
	Role  *string
	Title *string
}) *bool {

	userId, _, _, validToken := parseToken(args.Token)
	if !validToken {
		ret := false
		return &ret
	}

	ctx, collection := GetMongo("meetingItems")
	_, err := collection.InsertOne(
		ctx,
		bson.D{
			{"member", userId},
			{"date", args.Date},
			{"roleName", args.Role}, // In db we are now using roleName, since it's in the unique index.
			{"title", args.Title},
		},
	)
	if err != nil {
		log.Printf("Book error: %v. Input: %s, %s, %s\n", err, args.Date, *(args.Role), *(args.Title))
		ret := false
		return &ret
	}
	log.Printf("Successfully booked: %s, %s, %s\n", args.Date, *(args.Role), *(args.Title))
	ret := true
	return &ret
}
