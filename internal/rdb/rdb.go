package rdb

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/erfanshekari/go-talk/config"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

type RedisClient struct {
	Client *redis.Client
}

func (r *RedisClient) Close() {
	log.Println("Closing Redis Client...")
	err := r.Client.Close()
	if err != nil {
		log.Println(err)
	}
}

func (r *RedisClient) Ping() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	statusCmd := r.Client.Ping(ctx)
	if statusCmd != nil {
		response, err := statusCmd.Result()
		if err != nil {
			return false, err
		} else {
			if response == "PONG" {
				return true, nil
			}
		}
	}
	return false, errors.New("Somthing is Wrong with Redis Client!")
}

var singleInstance *RedisClient

func GetInstance(conf *config.Config) *RedisClient {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("**  Creating Redis Client...")

			if conf != nil && conf.Server.Debug {
				err := godotenv.Load()
				if err != nil {
					log.Fatal(err)
				}
			}
			options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
			if err != nil {
				log.Fatal(err)
			}
			singleInstance = &RedisClient{
				Client: redis.NewClient(options),
			}
		}
	}
	return singleInstance
}
