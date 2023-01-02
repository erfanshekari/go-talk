package models

import (
	"github.com/erfanshekari/go-talk/config"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollection(collectionName string, conf *config.Config) *mongo.Collection {
	client := mdbc.GetInstance(nil).Client
	return client.Database(conf.Database.Name).Collection(collectionName)
}
