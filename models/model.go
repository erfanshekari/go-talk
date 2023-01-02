package models

import (
	"context"

	"github.com/erfanshekari/go-talk/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	GetCollection() *mongo.Collection
	Migrate()
}

func Migrate(c context.Context, conf *config.Config) {
	MigrateUser(c, conf)
}
