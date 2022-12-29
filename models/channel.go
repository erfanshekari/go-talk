package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id,omitempity"`
	Name      string               `json:"name" bson:"name,omitempity"`
	Members   []primitive.ObjectID `json:"members" bson:"members,omitempity"`
	StartDate primitive.Timestamp  `json:"start_date" bson:"start_date,omitempity"`
}

func (c *Channel) GetMembers(ctx *context.Context) (*[]User, error) {
	// ...
	return nil, nil
}
