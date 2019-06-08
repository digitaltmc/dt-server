package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (_ *Resolver) Book(args struct {
	Token    string
	Date     string
	RoleName *string
	Title    *string
}) *bool {

	userId, _, _, validToken := parseToken(args.Token)
	if !validToken {
		ret := false
		return &ret
	}

	mi := MeetingItem{
		Date:     args.Date,
		RoleName: *args.RoleName,
		Member:   userId,
		Title:    *args.Title,
	}
	ctx, collection := GetMongo("meetingItems")
	_, err := collection.InsertOne(
		ctx,
		mi,
	)
	if err != nil {
		log.Printf("Book error: %v. Input: %s, %s, %s\n", err, args.Date, *(args.RoleName), *(args.Title))
		ret := false
		return &ret
	}
	log.Printf("Successfully booked: %s, %s, %s\n", args.Date, *(args.RoleName), *(args.Title))
	ret := true
	return &ret
}

func (_ *Resolver) Unbook(args struct {
	Token    string
	Date     string
	RoleName *string
}) *bool {

	_, _, _, validToken := parseToken(args.Token)
	if !validToken {
		ret := false
		return &ret
	}

	ctx, collection := GetMongo("meetingItems")
	_, err := collection.DeleteOne(
		ctx,
		bson.D{
			//			{"member", userId},
			{"date", args.Date},
			{"roleName", args.RoleName},
		},
	)
	if err != nil {
		log.Printf("Unbook error: %v. Input: %s, %s.\n", err, args.Date, *(args.RoleName))
		ret := false
		return &ret
	}
	log.Printf("Successfully unbooked: %s, %s.\n", args.Date, *(args.RoleName))
	ret := true
	return &ret
}
