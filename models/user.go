package models

// import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	// ID     primitive.ObjectID `json:"_id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
}
