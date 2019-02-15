package main

import (
  "fmt"
	"github.com/graph-gophers/graphql-go"

 	"github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/bson/primitive"
)

// 	"log"
// 	"net/http"
// 
// 	"github.com/graph-gophers/graphql-go/relay"

//  "github.com/friendsofgo/graphiql"
//  "github.com/mnmtanish/go-graphiql"


// 	"github.com/rs/cors"


type Resolver struct{}

//----------

type Meeting struct {
  id graphql.ID
  agenda []MeetingItem
  date string
}

type MeetingItem struct {
  role Role
  duration float64
  title string
}

type Role struct {
  name RolesEnum
  member Person
}

type Person struct {
  Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
  Name string
  Password string
  Mobile string
  Email string
  IsMember bool
  JoinedSince string
  MembershipUntil string
  Achievements []MeetingItem
}

type PersonInput struct {
  Name string
  Password string
}
//  Mobile string
//  Email string
//  IsMember bool
//  JoinedSince string
//  MembershipUntil string
//  Achievements []MeetingItem


type RolesEnum int
const (
  TMD RolesEnum = 0
  TTM RolesEnum = 1
  TTIE RolesEnum = 2
  GE RolesEnum = 3
  AhCounter RolesEnum = 4
  Grammarian RolesEnum = 5
  Timer RolesEnum = 6
  ShareMaster RolesEnum = 7
  Speaker RolesEnum = 8
  IE RolesEnum = 9
)

type OfficersEnum int
const (
  President OfficersEnum = 0
  VPE OfficersEnum = 1
  VPM OfficersEnum = 2
  VPPR OfficersEnum = 3
  Treasurer OfficersEnum = 4
  Secretary OfficersEnum = 5
  SAA OfficersEnum = 6
)

//---------- Query

func (_ *Resolver) Hello() string { return "Hello, world!" }


//----------

// type PersonResolver struct {
//   p *Person
// }

func (_ *Resolver) Register(arg *struct {Person *PersonInput}) *string {
  var succ = "true"
  var fail = "false"

  ctx, collection := GetMongo("person")
  cnt, err := collection.Count(
    ctx,
    bson.D{
      {"name", arg.Person.Name},
    },
  )
  if err != nil {
    fmt.Println(err)
  }
  if cnt != 0 {
    fmt.Printf("User already exists: %v", arg.Person.Name)
    return &fail
  }

  collection.InsertOne(
    ctx,
    bson.D{
      {"name", arg.Person.Name},
      {"password", arg.Person.Password},
    },
  )
  return &succ
}

func (_ *Resolver) Login(arg *struct {User, Password string}) *graphql.ID {
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

  var fail = graphql.ID("0")
  if err != nil {
    fmt.Println(err)
    return &fail
  }
  fmt.Println(p)
  var succ = graphql.ID(p.Id.Hex())
  return &succ
}

/*
mutation {
  register(person: {name: "owen", password: "pwd"})
}

{
  login(user: "owen", password: "pwd")
}

*/
