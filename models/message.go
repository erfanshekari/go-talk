package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageType string

const (
	Text MessageType = "text"
	File MessageType = "file"
)

type Message struct {
	Type       MessageType          `json:"type" bson:"type,omitempity"`
	From       primitive.ObjectID   `json:"from" bson:"from,omitempity"`
	Channel    primitive.ObjectID   `json:"channel" bson:"channel,omitempity"`
	Timestamp  primitive.Timestamp  `json:"timestamp" bson:"timestamp,omitempity"`
	Payload    interface{}          `json:"payload" bson:"payload,omitempity"`
	ReceivedBy []primitive.ObjectID `json:"delivered" bson:"delivered"`
	SeenBy     []primitive.ObjectID `json:"seen_by" bson:"seen_by"`
}
