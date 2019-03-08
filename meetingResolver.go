package main

import (
  "fmt"
  "regexp"
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
  agenda []*MeetingItem
  date string
}

type MeetingItem struct {
  Role Role
  Duration float64
  Title string
}

type Role struct {
  Name MeetingRolesEnum
  Member Person
}

type Person struct {
  Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
  Name string
  Password string
  Mobile string
  Email string
  IsMember bool
  OfficerRole OfficersEnum
  JoinedSince string
  MembershipUntil string
  Achievements []*MeetingItem
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


type MeetingRolesEnum string
const (
  TMD MeetingRolesEnum = "TMD"
  TTM MeetingRolesEnum = "TTM"
  TTIE MeetingRolesEnum = "TTIE"
  GE MeetingRolesEnum = "GE"
  AhCounter MeetingRolesEnum = "AhCounter"
  Grammarian MeetingRolesEnum = "Grammarian"
  Timer MeetingRolesEnum = "Timer"
  ShareMaster MeetingRolesEnum = "ShareMaster"
  Speaker MeetingRolesEnum = "Speaker"
  IE MeetingRolesEnum = "IE"
  RolePresident MeetingRolesEnum = "RolePresident"
  RoleSAA MeetingRolesEnum = "RoleSAA"
  RoleVPM MeetingRolesEnum = "RoleVPM"
  RoleVPE MeetingRolesEnum = "RoleVPE"
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


//---------- Mutations

func (_ *Resolver) Register(arg *struct {Person *PersonInput}) *graphql.ID {

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
    fmt.Printf("User already exists: %v\n", arg.Person.Name)
    return nil
  }

  insertRes, ins_err := collection.InsertOne(
    ctx,
    bson.D{
      {"name", arg.Person.Name},
      {"password", arg.Person.Password},
    },
  )
  if ins_err != nil {
    fmt.Printf("Insert error: %v\n", ins_err)
    return nil
  }

  // fmt.Printf("%T", insertRes.InsertedID) // InsertedID will be of type primitive.ObjectID
	re := regexp.MustCompile(`\"(.+?)\"`)
  ret := graphql.ID(re.FindStringSubmatch(insertRes.InsertedID.(primitive.ObjectID).String())[1])
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

// User can book a role if the role is not yet taken.
// User may book many roles in a meeting. - As long as they can handle.
func (_ *Resolver) Book(arg *struct {
  Person graphql.ID
  Date string
  Role MeetingRolesEnum
  Title *string
}) *MeetingResolver {
  ctx, collection := GetMongo("meeting")

  res := collection.FindOne(
    ctx,
    bson.D{
      {"date", arg.Date},
    },
  )
  fmt.Println(res)
  m := Meeting{graphql.ID(0), nil, ""}
  mr := &MeetingResolver{&m}
  return mr
//  if err != nil {
//    fmt.Println(err)
//  }
//  if cnt != 0 {
//    fmt.Printf("User already exists: %v", arg.Person.Name)
//    return &fail
//  }
//
//  collection.InsertOne(
//    ctx,
//    bson.D{
//      {"name", arg.Person.Name},
//      {"password", arg.Person.Password},
//    },
//  )
}

//---------- MeetingResolver

type MeetingResolver struct {
  m *Meeting
}

func (r *MeetingResolver) Id() graphql.ID {
	return r.m.id
}

func (r *MeetingResolver) Agenda() *[]*MeetingItemResolver {
  ret := make([]*MeetingItemResolver, len(r.m.agenda))
  for i, v := range r.m.agenda {
    ret[i] = &MeetingItemResolver{v}
  }
	return &ret
}

func (r *MeetingResolver) Date() *string {
	return &r.m.date
}


//---------- MeetingItemResolver

type MeetingItemResolver struct {
  m *MeetingItem
}

func (r *MeetingItemResolver) Role() *RoleResolver {
  rr := RoleResolver{&r.m.Role}
	return &rr
}

func (r *MeetingItemResolver) Duration() *float64 {
	return &r.m.Duration
}

func (r *MeetingItemResolver) Title() *string {
	return &r.m.Title
}

//---------- RoleResolver

type RoleResolver struct {
  r *Role
}

func (r *RoleResolver) Name() *MeetingRolesEnum {
	return &r.r.Name
}

func (r *RoleResolver) Member() *PersonResolver {
  rr := PersonResolver{&r.r.Member}
	return &rr
}

//---------- PersonResolver

type PersonResolver struct {
  r *Person
}

func (r *PersonResolver) Id() graphql.ID {
  ret := graphql.ID(r.r.Id.String())
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

func (r *PersonResolver) Email() *string {
	return &r.r.Email
}

func (r *PersonResolver) OfficerRole() *OfficersEnum {
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

func (r *PersonResolver) Achievements() *[]*MeetingItemResolver {
	ret := makeMeetingItemResolver(r.r.Achievements)
	return &ret
}

