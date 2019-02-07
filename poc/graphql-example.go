// Embedded in this article https://medium.com/p/c98e491015b6
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/rs/cors"
)

//////// MAIN ////////
func main() {
	// Create a handler for /graphql which passes cors for remote requests
	http.Handle("/graphql", cors.Default().Handler(&relay.Handler{Schema: graphqlSchema}))

	// ListenAndServe starts an HTTP server with a given address and handler.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//////// MONGODB ////////
const host string = "localhost"

// Cleanup will remove all mock data from the database.
func Cleanup(col string) {
	log.Println("Cleaning up MongoDB...")
	ctx, collection := GetMongo(col)

	_, err := collection.DeleteMany(ctx,
		bson.D{})
	if err != nil {
		log.Fatal(err)
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
	var options options.ClientOptions
	maxWait := time.Duration(5 * time.Second)
	options.SetConnectTimeout(maxWait)
//	options.SetAuth(auth)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, "mongodb://"+host+":27017", &options)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("minimalGraphql").Collection(col)
	return ctx, collection
}

//////// GRAPHQL ////////
var graphqlSchema *graphql.Schema

// Schema describes the data that we ask for
var Schema = `
    schema {
        query: Query
    }
    # The Query type represents all of the entry points.
    type Query {
        post(slug: String!): Post
    }
    type Post {
        id: ID!
        slug: String!
        title: String!
    }
    `

//////// INIT ////////
func init() {
	// MustParseSchema parses a GraphQL schema and attaches the given root resolver.
	// It returns an error if the Go type signature of the resolvers does not match the schema.
	graphqlSchema = graphql.MustParseSchema(Schema, &Resolver{})

	log.Println("Seeding mock data to MongoDB")
	// Call GetMongo, session and reference to the post collection
	ctx, collection := GetMongo("post")
	// Close the session so its resources may be put back in the pool or collected, depending on the case.

	// Cleanup finds all documents matching the provided selector document
	// and removes them from the database. So we make sure the db is empty before inserting mock data.
	Cleanup("post")

	// The mock data that we insert.
	_, err := collection.InsertMany(
		ctx,
		[]interface{}{
			bson.D{
				{"ID", 1},
				{"title", "First post"},
				{"slug", "first-post"},
			},
			bson.D{
				{"ID", 2},
				{"title", "Second post"},
				{"slug", "second-post"},
			},
			bson.D{
				{"ID", 3},
				{"title", "Third post"},
				{"slug", "third-post"},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mock data added successfully!")
}

//////// RESOLVER ////////
// In order to respond to queries, a schema needs to have resolve functions for all fields.
// Go’s structs are typed collections of fields. They’re useful for grouping data together to form records.
type Resolver struct{}

type post struct {
	ID    graphql.ID
	Slug  string
	Title string
}

type postResolver struct {
	s *post
}

type searchResultResolver struct {
	result interface{}
}

// Slices can be created with the built-in make function; this is how we create dynamically-sized arrays.
var postData = make(map[string]*post)

// Post resolves the Post queries.
func (r *Resolver) Post(args struct{ Slug string }) *postResolver {
	// One result is a pointer to type post.
	oneResult := &post{}

	// Call GetMongo, session and reference to the post collection
	ctx, collection := GetMongo("post")
	// Close the session so its resources may be put back in the pool or collected, depending on the case.

	// Inside the collection, find by slug and return all fields.
	//err := collection.Find(bson.M{"slug": args.Slug}).Select(bson.M{}).One(&oneResult)
	cur, err := collection.Find(
		ctx,
		bson.D{
      {"slug", args.Slug},
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		cur.Decode(oneResult)
	}

	// Make a type postResolver out of oneResult.
	if s := oneResult; s != nil {
		return &postResolver{oneResult}
	}
	return nil
}

// Resolve each field to respond to queries.
func (r *postResolver) ID() graphql.ID {
	return r.s.ID
}

func (r *postResolver) Slug() string {
	return r.s.Slug
}

func (r *postResolver) Title() string {
	return r.s.Title
}
