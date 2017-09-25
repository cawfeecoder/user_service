package redis

import (
	"log"

	"github.com/go-redis/redis"
)

var c *redis.Client
var d *redis.Client

func init() {

	c = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	d = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	log.Printf("Successfully connected to the Database!")
}

//GetOriginClient - Get a connection to Origin Token Database
func GetOriginClient() *redis.Client {
	return c
}

//GetDerivClient - Get a connection to Derivative Token Database
func GetDerivClient() *redis.Client {
	return d
}
