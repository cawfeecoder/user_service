package main

import (
	"log"

	pb "github.com/nfrush/user_service/user_service_pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "testClient"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect due to: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	createAccountRequest := &pb.CreateAccountRequest{
		FirstName:     "Test",
		MiddleInitial: "G.",
		LastName:      "Test",
		Username:      "test",
		Password:      "test123",
		Email:         "test@test.com",
		PhoneNumber:   "4439740421",
	}

	r, err := c.CreateAccount(context.Background(), createAccountRequest)
	if err != nil {
		log.Fatalf("Could not create account: %v", err)
	}
	log.Printf("Received a message of Status %d with Message %s", r.StatusType, r.StatusMessage)
}
