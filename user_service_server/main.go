package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/nfrush/user_service/user_service_pb/user_service"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return nil, &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message", nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Server(lis); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
	log.Printf("Started server successfully on port %s", port);
}
