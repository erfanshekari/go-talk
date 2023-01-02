package models

import (
	"context"
	"log"
	"time"

	"github.com/erfanshekari/go-talk/config"
	// ctx "github.com/erfanshekari/go-talk/context"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init Model Name
const modelName string = "user"

type UserStatus string

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
	Banned  UserStatus = "banned"
)

type UserDetails struct {
	ID       primitive.ObjectID   `json:"_id" bson:"_id,omitempity"`
	UserID   string               `json:"user_id" bson:"user_id,omitempity"`
	Status   UserStatus           `json:"status" bson:"status,omitempity"`
	LastSeen primitive.Timestamp  `json:"lastseen" bson:"lastseen,omitempity"`
	Contacts []primitive.ObjectID `json:"contacts" bson:"contacts"`
	Channels []primitive.ObjectID `json:"channels" bson:"channels"`
}

type User struct {
	UserID  string
	Details *UserDetails
}

func MigrateUser(c context.Context, conf *config.Config) {
	db := mdbc.GetInstance(nil).Client.Database(conf.Database.Name)
	err := db.CreateCollection(c, modelName)
	if err != nil {
		log.Println(err)
	}
	coll := db.Collection(modelName)
	indexUnique := true
	_, err = coll.Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{primitive.E{
			Key:   "user_id",
			Value: 1,
		}},
		Options: &options.IndexOptions{
			Unique: &indexUnique,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func (u *User) GetCollection() *mongo.Collection {
	return GetCollection(modelName)
}

func (u *User) newUser(userID string, c context.Context) (*UserDetails, error) {
	coll := u.GetCollection()
	details := UserDetails{
		ID:       primitive.NewObjectID(),
		UserID:   userID,
		LastSeen: primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0},
		Status:   Offline,
		Contacts: []primitive.ObjectID{},
		Channels: []primitive.ObjectID{},
	}
	log.Println(details.LastSeen, details.LastSeen.T, details.LastSeen.I)
	_, err := coll.InsertOne(c, details)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &details, nil
}

func (u *User) GetDetails(c context.Context) *UserDetails {
	if u.Details == nil {
		coll := u.GetCollection()
		result := coll.FindOne(c, bson.D{primitive.E{
			Key:   "user_id",
			Value: u.UserID,
		}})
		if err := result.Err(); err != nil {
			// insert user to db
			if err == mongo.ErrNoDocuments {
				d, err := u.newUser(u.UserID, c)
				if err != nil {
					log.Println(err.Error())
				}
				u.Details = d
			} else {
				log.Println(err.Error())
			}
		} else {
			decoded, err := result.DecodeBytes()
			if err != nil {
				log.Println(err.Error())
			} else {
				var dtls UserDetails
				err = bson.Unmarshal(decoded, &dtls)
				if err != nil {
					log.Println(err.Error())
				} else {
					u.Details = &dtls
				}
			}
		}
	}
	return u.Details
}

func (u *User) UpdateActivity(ctx *context.Context) error {
	// ... Update LastSeen To Now
	return nil
}

func (u *User) SetStatus(ctx *context.Context, status UserStatus) error {
	// ...
	return nil
}

func (u *User) NotifyContacts(ctx *context.Context) error {
	// ...
	return nil
}
