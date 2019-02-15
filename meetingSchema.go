package main

// import (
// 	"log"
// 	"net/http"
// 
// 	"github.com/graph-gophers/graphql-go"
// 	"github.com/graph-gophers/graphql-go/relay"
// )

// 	"github.com/rs/cors"

// 	"github.com/mongodb/mongo-go-driver/bson"

//  "github.com/friendsofgo/graphiql"
//  "github.com/mnmtanish/go-graphiql"


// Schema describes the data that we ask for
var Schema = `
    schema {
      query: Query
      mutation: Mutation
    }
    # The Query type represents all of the entry points.
    type Query {
      hello: String!

      login(user: String!, password: String!): ID
#      meetings: [Meeting]
#      meeting(date: String): Meeting
    }

    type Mutation {
      register(person: PersonInput): String
#      book(date: String!, role: RolesEnum!, title: String): Meeting
    }

    type Post {
        id: ID!
        slug: String!
        title: String!
    }

    type Meeting {
      id: ID!
      agenda: [MeetingItem]
      date: String
    }

    type MeetingItem {
      role: Role
      duration: Float
      title: String
    }

    type Role {
      name: RolesEnum
      member: Person
    }

    type Person {
      id: ID!
      name: String
      mobile: String
      email: String
      password: String
      isMember: Boolean
      joinedSince: String
      membershipUntil: String
      achievements: [ MeetingItem ]
    }
    input PersonInput {
      name: String!
      password: String!
#      mobile: String
#      email: String
#      isMember: Boolean!
#      joinedSince: String
#      membershipUntil: String
#      achievements: [ MeetingItem ]
    }

    enum RolesEnum {
      TMD
      TTM
      TTIE
      GE
      AhCounter
      Grammarian
      Timer
      ShareMaster
      Speaker
      IE
    }

    enum OfficersEnum {
      President
      VPE
      VPM
      VPPR
      Treasurer
      Secretary
      SAA
    }
    `

