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

func Migrate(c context.Context, conf *config.ConfigAtrs) {
	var (
		user User
	)
	user.Migrate(c, conf)
}
