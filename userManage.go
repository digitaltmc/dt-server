package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (_ *Resolver) Register(arg *struct{ Person *PersonInput }) *string {
	if alreadyExists, err := arg.Person.Exists(); alreadyExists == true || err != nil {
		log.Printf("User already exists: %v\n", arg.Person.Name)
		return nil
	}

	ctx, collection := GetMongo("person")
	insertRes, err := collection.InsertOne(
		ctx,
		bson.D{
			{"name", arg.Person.Name},
			{"password", arg.Person.Password},
			{"email", arg.Person.Email},
			{"mobile", CheckNilString(arg.Person.Mobile)},
		},
	)
	if err != nil {
		log.Printf("Insert error: %v\n", err)
		return nil
	}

	ret, err := createToken(insertRes.InsertedID.(primitive.ObjectID).Hex(), arg.Person.Name, arg.Person.Email)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &ret
}

func (_ *Resolver) WxLogin(arg *struct{ Code string }) string {
	wxInfo, err := getwxLoginResult(arg.Code)
	if wxInfo.Openid != "" {
		ctx, collection := GetMongo("person")
		c := collection.FindOne(
			ctx,
			bson.D{
				{"openid", wxInfo.Openid},
			},
		)
		var p Person
		var err = c.Decode(&p)

		if err != nil {
			log.Println(err)
			openid := wxInfo.Openid
			return openid
		}
		log.Println(p)
		id := p.Id.Hex()
		return id
	} else {
		e := err.Error()
		return e
	}
}

func (_ *Resolver) Login(arg *struct{ User, Password string }) *string {
	log.Println("Login: ", arg)
	ctx, collection := GetMongo("person")
	c := collection.FindOne(
		ctx,
		bson.D{
			{"name", arg.User},
			{"password", arg.Password},
		},
	)

	var p Person
	var err = c.Decode(&p)

	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println(p)
	log.Println(p.Id.Hex())

	ret, err := createToken(p.Id.Hex(), p.Name, p.Email)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &ret
}

func (_ *Resolver) DeleteUser(arg *struct{ Username string }) *string {
	log.Println("DeleteUser:", arg)
	ctx, collection := GetMongo("person")
	_, err := collection.DeleteOne(
		ctx,
		bson.D{
			{"name", arg.Username},
		},
	)
	if err != nil {
		log.Println("Fail to delete user: ", err)
		res := "Fail to delete user: " + arg.Username
		return &res
	}
	res := "Successfully delete user: " + arg.Username
	return &res
}
