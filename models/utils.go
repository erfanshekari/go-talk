package models

import (
	"log"

	"github.com/erfanshekari/go-talk/config"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollection(collectionName string, conf *config.ConfigAtrs) *mongo.Collection {
	log.Println(conf.DatabaseName)
	client := mdbc.GetInstance(nil).Client
	return client.Database(conf.DatabaseName).Collection(collectionName)
}
