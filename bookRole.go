package main

import (
	"github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson"
)

func (_ *Resolver) Book(args struct {
	Token string
	Date  graphql.Time
	Role  *string
	Title *string
}) bool {

	userId, _, _, validToken := parseToken(args.Token)
	if !validToken {
		return false
	}

	ctx, collection := GetMongo("meeting")
	_, error := collection.InsertOne(
		ctx,
		bson.D{
			{"date", args.Date},
			{"role", args.Role},
			{"member", userId},
			{"title", args.Title},
		},
	)
}
