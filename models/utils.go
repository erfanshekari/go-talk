package models

import (
	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollection(collectionName string) *mongo.Collection {
	client := mdbc.GetInstance(nil).Client
	return client.Database(global.GetInstance(nil).Config.Database.Name).Collection(collectionName)
}
