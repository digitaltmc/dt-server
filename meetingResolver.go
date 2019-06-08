package main

import (
	"context"
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
	Date   string         `bson:"date"`
	Agenda []*MeetingItem `bson:"agenda"`
}
type MeetingItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Date     string             `bson:"date"`
	RoleName string             `bson:"roleName"`
	Member   string             `bson:"member"`
	Title    string             `bson:"title"`
	Duration string             `bson:"duration"`
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

//----------

func GetMeeting(currentdate string) *Meeting {
	var meetingItems []*MeetingItem
	ctx, collection := GetMongo("meetingItems")
	filter := bson.D{{"date", currentdate}}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Fail to get meeting: ", err)
		return nil
	}

	for cur.Next(ctx) {
		var elem MeetingItem
		err := cur.Decode(&elem)
		log.Println(elem)
		if err != nil {
			log.Fatal(err)
		}
		meetingItems = append(meetingItems, &elem)
	}

	meeting := Meeting{
		Date:   currentdate,
		Agenda: meetingItems,
	}
	return &meeting
}
func (_ *Resolver) Meeting(args struct {
	Date string
}) *meetingResolver {
	currentdate := args.Date
	log.Println("[current date] ", currentdate)
	meeting := GetMeeting(currentdate)
	if meeting != nil {
		return &meetingResolver{meeting}
	}
	return nil
}

//----------

func DecodeBookList(cursor *mongo.Cursor) []Meeting {
	// Largest size is 52 weeks
	meetings := make([]Meeting, 52)
	i := 0
	for cursor.Next(context.Background()) {
		cursor.Decode(&meetings[i])
		i++
	}
	log.Println("result length", i)
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
	var currentdate = time.Now()
	limiteddate := currentdate.AddDate(-1, 0, 0)
	log.Println("!!!!!!", limiteddate)
	//TODO set the limitation of query date
	// filter := bson.D{{"date", bson.D{{"$lt", currentdate},{"$gt", limiteddate},}},}
	filter := bson.D{{"date", bson.D{{"$lt", currentdate}}}}
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		log.Println("!!!!!!!!!!!!", err)
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

type meetingResolver struct {
	m *Meeting
}

// func (r *meetingResolver) ID() graphql.ID { return graphql.ID(r.m.ID.Hex()) }

// Please convert the time to UTC as the graphql.Time belong to ISO with time.RFC3339
// func (r *meetingResolver) Date() graphql.Time {
func (r *meetingResolver) Date() string {
	// log.Println("The UTC Time ", r.m.Date.Time.UTC())
	return r.m.Date
}
func (r *meetingResolver) Agenda() *[]*meetingItemResolver {
	l := make([]*meetingItemResolver, len(r.m.Agenda))

	for i, v := range r.m.Agenda {
		log.Println(v)
		l[i] = &meetingItemResolver{v}
		log.Println("Agenda RoleName:", v.RoleName)
	}
	return &l
}

type meetingItemResolver struct {
	mi *MeetingItem
}

func (r *meetingItemResolver) ID() graphql.ID { return graphql.ID(r.mi.ID.Hex()) }
func (r *meetingItemResolver) Date() *string {
	return &r.mi.Date
}
func (r *meetingItemResolver) RoleName() *string {
	return &r.mi.RoleName
}
func (r *meetingItemResolver) Member() *PersonResolver {
	ctx, collection := GetMongo("person")
	// log.Println("member id: ", r.mi.Member.Hex())
	mid, err := primitive.ObjectIDFromHex(r.mi.Member)
	if err != nil {
		log.Fatal("Fail to convert ObjectId:", r.mi.Member)
	}
	log.Println("Find id: ", mid)
	filter := bson.D{{"_id", mid}}

	var person Person
	err = collection.FindOne(ctx, filter).Decode(&person)
	if err != nil {
		log.Println("Fail to find member: ", err)
		return nil
	}
	return &PersonResolver{&person}
}

func (r *meetingItemResolver) Duration() *string { return &r.mi.Duration }
func (r *meetingItemResolver) Title() *string    { return &r.mi.Title }

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
