package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageType string

const (
	Text MessageType = "text"
	File MessageType = "file"
)

type Message struct {
	Channel primitive.ObjectID `json:"channel" bson:"channel,omitempity"`
	From    primitive.ObjectID `json:"from" bson:"from,omitempity"`
	Payload string             `json:"payload" bson:"payload,omitempity"`
	Type    MessageType        `json:"type" bson:"type,omitempity"`
}
