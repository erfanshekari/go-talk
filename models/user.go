package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID  string
	Details *UserDetails
}

type UserStatus string

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
	Deleted UserStatus = "deleted"
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

func (u *User) GetDetails(ctx *context.Context) *UserDetails {
	if u.Details == nil {
		// ... Get User from MongoDB
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
