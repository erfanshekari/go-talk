package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID      primitive.ObjectID   `json:"_id" bson:"_id,omitempity"`
	Name    string               `json:"name" bson:"name,omitempity"`
	Members []primitive.ObjectID `json:"members" bson:"members,omitempity"`
}
