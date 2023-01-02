package mdbc

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/erfanshekari/go-talk/config"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock = &sync.Mutex{}

type MongoClient struct {
	Client        *mongo.Client
	Context       *context.Context
	CancelContext *context.CancelFunc
}

func (m *MongoClient) Ping() (bool, error) {
	dbs, err := m.Client.ListDatabaseNames(*m.Context, bson.D{})
	if err != nil {
		return false, err
	} else {
		if len(dbs) >= 1 {
			return true, nil
		}
	}
	return true, errors.New("Somethig is wrong with MongoDB client")
}

func (m *MongoClient) Close() {
	log.Println("Closing MongoDB Client...")
	defer (*m.CancelContext)()
	err := m.Client.Disconnect(*m.Context)
	if err != nil {
		log.Println(err)
	}
}

var singleInstance *MongoClient

func GetInstance(conf *config.Config) *MongoClient {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("** Creating MongoDB Client...")

			if conf != nil && conf.Server.Debug {
				err := godotenv.Load()
				if err != nil {
					log.Fatal(err)
				}
			}

			uri := os.Getenv("MONGO_DB_URI")
			client, err := mongo.NewClient(options.Client().ApplyURI(uri))
			if err != nil {
				log.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			err = client.Connect(ctx)
			if err != nil {
				log.Fatal(err)
			}
			singleInstance = &MongoClient{
				Client:        client,
				Context:       &ctx,
				CancelContext: &cancel,
			}
		}
	}
	return singleInstance
}
