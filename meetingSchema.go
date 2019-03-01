package main

// import (
// 	"log"
// 	"net/http"
// 
// 	"github.com/graph-gophers/graphql-go"
// 	"github.com/graph-gophers/graphql-go/relay"
// )


// Schema describes the data that we ask for
var Schema = `
    schema {
      query: Query
      mutation: Mutation
    }
    # The Query type represents all of the entry points.
    type Query {
      hello: String!
      wxLogin(code: String!): String!

      login(user: String!, password: String!): ID
#      meetings: [Meeting]
#      meeting(date: String): Meeting
    }

    type Mutation {
      register(person: PersonInput): String
      book(date: String!, role: RolesEnum!, title: String): Meeting
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
      name: MeetingRolesEnum
      member: Person
    }

    type Person {
      id: ID!
      name: String
      mobile: String
      email: String
      password: String
      officerRole: OfficersEnum
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

    enum MeetingRolesEnum {
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
      RolePresident
      RoleSAA
      RoleVPM
      RoleVPE
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

