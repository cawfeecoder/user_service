package main

import (
	"log"
	"net"

	pb "github.com/nfrush/user_service/user_service_pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}

func (s *server) UpdateInfo(ctx context.Context, request *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {
	return &pb.UpdateInfoResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) DeleteAccount(ctx context.Context, request *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	return &pb.DeleteAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "This is a test message", OriginToken: "test", DerivativeToken: "test"}, nil
}
func (s *server) RefreshSession(ctx context.Context, request *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	return &pb.RefreshSessionResponse{StatusType: 1, StatusMessage: "This is a test message", DerivativeToken: "test"}, nil
}
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
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
	log.Printf("Started server successfully on port %s", port)
}
