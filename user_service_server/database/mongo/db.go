package db

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var c *mgo.Database

//GetSession - Get a session from the DB
func GetSession() *mgo.Database {
	session, err := mgo.Dial("localhost:27018")
	if err != nil {
		log.Fatalf("Error creating DB Connection: %v", err)
	}
	c = session.DB("test")
	return c
}
