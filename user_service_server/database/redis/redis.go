package redis

import (
	"log"

	"github.com/go-redis/redis"
)

var c *redis.Client

func init() {

	c = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Printf("Successfully connected to the Database!")
}

//GetClient - Get a connection from the DB
func GetClient() *redis.Client {
	return c
}
