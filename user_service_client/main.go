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

	createAccountRequest := &pb.LogoutRequest{
		//FirstName:     "Test",
		//MiddleInitial: "G.",
		//LastName:      "Test",
		Username:        "test",
		OriginToken:     "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0IiwiZXhwIjoxNTA2Mzg3Njg0LCJpYXQiOjE1MDYxMjg0ODQsImlzcyI6IkZydXNoIERldmVsb3BtZW50IExURCIsImp0aSI6Imh0dHA6Ly9leGFtcGxlLmNvbSJ9.HvL86CLLtOR_e6okWfxw5JBcuHhHumHZyFBAdpcghDxLQ5HyWc-MpWLbJG_ai9IUar0TypjTt5shkAJsss3Zfg",
		DerivativeToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJleUpoYkdjaU9pSklVelV4TWlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKaGRXUWlPaUowWlhOMElpd2laWGh3SWpveE5UQTJNemczTmpnMExDSnBZWFFpT2pFMU1EWXhNamcwT0RRc0ltbHpjeUk2SWtaeWRYTm9JRVJsZG1Wc2IzQnRaVzUwSUV4VVJDSXNJbXAwYVNJNkltaDBkSEE2THk5bGVHRnRjR3hsTG1OdmJTSjkuSHZMODZDTEx0T1JfZTZva1dmeHc1SkJjdUhoSHVtSFp5RkJBZHBjZ2hEeExRNUh5V2MtTXBXTGJKR19haTlJVWFyMFR5cGpUdDVzaGtBSnNzczNaZmciLCJleHAiOjE1MDYzMDExODIsImlhdCI6MTUwNjI5NzU4MiwiaXNzIjoiRnJ1c2ggRGV2ZWxvcG1lbnQgTFREIiwianRpIjoiaHR0cDovL2V4YW1wbGUuY29tIiwicm9sZXMiOlsiVXNlciJdfQ.hURNanQphN9kErE_7msP8wxOnDJPXSMmMTcHnEq-BZj486jtZlD_g5cP_EU-GWyWa7UIBvsJ0Kj65nCJc0OV5Q",
		//Password: "test123",
		//Email:         "test@test.com",
		//PhoneNumber:   "4439740421",
	}

	r, err := c.Logout(context.Background(), createAccountRequest)
	if err != nil {
		log.Fatalf("Could not create account: %v", err)
	}
	log.Printf("Received a message of Status %d with Message %s", r.StatusType, r.StatusMessage)
}
