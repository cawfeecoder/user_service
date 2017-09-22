package db

import (
	"crypto/tls"
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var c driver.Client

func init() {
	var err error

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	if err != nil {
		log.Fatalf("Error creating HTTP Connection: %v", err)
	}

	c, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("user_microservice_admin", "gqCuSI9p4q2"),
	})
	if err != nil {
		log.Fatalf("Error connecting to Database: %v", err)
	} else {
		log.Printf("Successfully connected to the Database!")
	}

}

//GetClient - Get a connection from the DB
func GetClient() driver.Client {
	return c
}
