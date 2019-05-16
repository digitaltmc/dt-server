package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/graph-gophers/graphql-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct{}

//----------

type Meeting struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Date   graphql.Time       `bson: "date"`
	Agenda []*MeetingItem     `bson: "agenda"`
}
type MeetingItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Role     string             `bson: "role"`
	Member   primitive.ObjectID `bson: "member"`
	Duration int                `bson: "duration"` // In seconds
	Title    string             `bson: "title"`
}

type Person struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string
	Password        string
	Mobile          string
	Email           string
	IsMember        bool
	OfficerRole     string
	JoinedSince     string
	MembershipUntil string
	Achievements    []*MeetingItem
}

type PersonInput struct {
	Name            string
	Password        string
	Email           string
	Mobile          *string `bson:",omitempty"`
	IsMember        *bool   `bson:",omitempty"`
	JoinedSince     *string `bson:",omitempty"`
	MembershipUntil *string `bson:",omitempty"`
}

func (p *PersonInput) Exists() (bool, error) {
	cnt, err := countItems("name", p.Name)
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}

	cnt, err = countItems("email", p.Email)
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

var currentID primitive.ObjectID

func (_ *Resolver) Hello() string { return "Hello, world!" }

//---------- Mutations

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

	ret, err := createToken(insertRes.InsertedID.(string), arg.Person.Name, arg.Person.Email)
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
			fmt.Println(err)
			openid := wxInfo.Openid
			return openid
		}
		fmt.Println(p)
		id := p.Id.Hex()
		return id
	} else {
		e := err.Error()
		return e
	}
}

func (_ *Resolver) Login(arg *struct{ User, Password string }) *string {
	fmt.Println(arg)
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
		fmt.Println(err)
		return nil
	}

	log.Println(p)

	ret, err := createToken(p.Id.String(), p.Name, p.Email)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &ret
}

func GetMeeting(currentdate graphql.Time) *Meeting {
	var meeting Meeting
	ctx, collection := GetMongo("meeting")
	filter := bson.D{{"date", currentdate}}
	err := collection.FindOne(ctx, filter).Decode(&meeting)
	if err != nil {
		fmt.Println("!!!!!!!!!!!!", err)
		return nil
	}
	return &meeting
}
func (_ *Resolver) Meeting(args struct {
	Date graphql.Time
}) *meetingResolver {
	currentdate := args.Date
	fmt.Println("[current date] ", currentdate)
	booker := GetMeeting(currentdate)
	if booker != nil {
		return &meetingResolver{booker}
	}
	return nil
}

func DecodeBookList(cursor *mongo.Cursor) []Meeting {
	// Largest size is 52 weeks
	meetings := make([]Meeting, 52)
	i := 0
	for cursor.Next(context.Background()) {
		cursor.Decode(&meetings[i])
		i++
	}
	fmt.Println("result length", i)
	return meetings[0:i]
}
func ContainsKey(doc bson.Raw, key ...string) bool {
	_, err := doc.LookupErr(key...)
	if err != nil {
		return false
	}
	return true
}
func GetBookers() []Meeting {
	ctx, collection := GetMongo("meeting")
	var currentdate = graphql.Time{Time: time.Now()}
	limiteddate := currentdate.AddDate(-1, 0, 0)
	fmt.Println("!!!!!!", limiteddate)
	//TODO set the limitation of query date
	// filter := bson.D{{"date", bson.D{{"$lt", currentdate},{"$gt", limiteddate},}},}
	filter := bson.D{{"date", bson.D{{"$lt", currentdate}}}}
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		fmt.Println("!!!!!!!!!!!!", err)
		return nil
	}
	return DecodeBookList(cursor)
}

func (_ *Resolver) Meetings() *[]*meetingResolver {
	bookers := GetBookers()
	if bookers != nil {

		bookerlist := make([]*meetingResolver, len(bookers))

		for i, _ := range bookers {
			bookerlist[i] = &meetingResolver{&bookers[i]}
		}

		return &bookerlist
	}
	return nil
}

func (_ *Resolver) BookOld(args struct {
	Token string
	Date  graphql.Time
	Role  *string
	Title *string
}) *meetingResolver {

	userId, _, _, validToken := parseToken(args.Token)
	if !validToken {
		return nil
	}

	currentdate := args.Date
	booker := GetMeeting(currentdate)

	if booker == nil {
		ctx, collection := GetMongo("meeting")
		_, error := collection.InsertOne(
			ctx,
			bson.D{
				{"date", currentdate},
				{"agenda",
					bson.A{
						bson.D{
							{"role", args.Role},
							{"duration", 7.30},
							{"title", args.Title},
							{"member", userId},
						},
					},
				},
			},
		)
		fmt.Println("Inserted ", error)
		booker = GetMeeting(currentdate)
		// ToDo insert new meeting
	}
	// fmt.Println(meeting,meeting.id)
	return &meetingResolver{booker}
}

type meetingResolver struct {
	m *Meeting
}

func (r *meetingResolver) ID() graphql.ID { return graphql.ID(r.m.ID.Hex()) }

// Please convert the time to UTC as the graphql.Time belong to ISO with time.RFC3339
func (r *meetingResolver) Date() graphql.Time {
	fmt.Println("Graphql Time ", r.m.Date)
	fmt.Println("The UTC Time ", r.m.Date.Time.UTC())
	return r.m.Date
}
func (r *meetingResolver) Agenda() *[]*meetingItemResolver {
	l := make([]*meetingItemResolver, len(r.m.Agenda))

	for i, _ := range r.m.Agenda {
		l[i] = &meetingItemResolver{r.m.Agenda[i]}
		fmt.Println("Agenda RoleName", r.m.Agenda[i].Role)
	}
	return &l
}

type meetingItemResolver struct {
	mi *MeetingItem
}

func (r *meetingItemResolver) ID() graphql.ID { return graphql.ID(r.mi.ID.Hex()) }
func (r *meetingItemResolver) Role() *string {
	return &r.mi.Role
}
func (r *meetingItemResolver) Member() *PersonResolver {
	var person Person
	ctx, collection := GetMongo("person")
	filter := bson.D{{"_id", r.mi.Member}}
	err := collection.FindOne(ctx, filter).Decode(&person)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &PersonResolver{&person}
}

func (r *meetingItemResolver) Duration() *int { return &r.mi.Duration }
func (r *meetingItemResolver) Title() *string { return &r.mi.Title }

//---------- PersonResolver

type PersonResolver struct {
	r *Person
}

func (r *PersonResolver) Id() graphql.ID {
	ret := graphql.ID(r.r.Id.Hex())
	return ret
}

func (r *PersonResolver) Name() *string {
	return &r.r.Name
}

func (r *PersonResolver) Password() *string {
	return &r.r.Password
}

func (r *PersonResolver) Mobile() *string {
	return &r.r.Mobile
}

func (r *PersonResolver) Email() string {
	return r.r.Email
}

func (r *PersonResolver) OfficerRole() *string {
	return &r.r.OfficerRole
}

func (r *PersonResolver) IsMember() *bool {
	return &r.r.IsMember
}

func (r *PersonResolver) JoinedSince() *string {
	return &r.r.JoinedSince
}

func (r *PersonResolver) MembershipUntil() *string {
	return &r.r.MembershipUntil
}

func (r *PersonResolver) Achievements() *[]*meetingItemResolver {
	ret := makeMeetingItemResolver(r.r.Achievements)
	return &ret
}

/*
1. Data structor: enum/meetingItem & Role
2. book's paramter
3. logic to get the agenda from book
mutation {
  register(person:{name:"Wow",password:"world",email:"aaa@sap.com",mobile:"888888"})
}
{
  login(user:"Wow",password:"world")
}

mutation {
  book(date: "2019-03-11T00:00:00Z", role: Speaker, title: "Hey buddy") {
    date
    agenda {
      role
      title
      duration
      member{
        name
        mobile
        email
        id
      }
    }
  }
}
*/
