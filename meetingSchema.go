package main

// import (
// 	"log"
// 	"net/http"
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

      login(user: String!, password: String!): String
      meetings: [Meeting]
      meeting(date: String!): Meeting!
    }

    type Mutation {
      register(person: PersonInput): String
      book(token: String!, date: String!, roleName: String, title: String): Boolean
    }

    scalar Time

    type Post {
        id: ID!
        slug: String!
        title: String!
    }

    type Meeting {
      id: ID!
      agenda: [MeetingItem]
      date: String!
    }

    type MeetingItem {
      id: ID!
      role: MeetingRolesEnum
      member: Person
      duration: String
      title: String
    }

    type Person {
      id: ID!
      name: String
      mobile: String
      email: String!
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
      mobile: String
      email: String!
      isMember: Boolean
      joinedSince: String
      membershipUntil: String
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
      President
      SAA
      VPM
      VPE
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
