package main

import (
	"log"
	"net"

	driver "github.com/arangodb/go-driver"
	pb "github.com/nfrush/user_service/user_service_pb"
	db "github.com/nfrush/user_service/user_service_server/database"
	"github.com/nfrush/user_service/user_service_server/models/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	log.Printf("Creating Account For %s, %s %s with username %s and password %s using email %s and phone number %s", request.LastName, request.MiddleInitial, request.FirstName, request.Username, request.Password, request.Email, request.PhoneNumber)
	generatehash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	doc := modelUser.User{
		FirstName:     request.Username,
		MiddleInitial: request.MiddleInitial,
		LastName:      request.LastName,
		Username:      request.Username,
		Password:      string(generatehash),
		Status:        1,
		Email:         request.Email,
		PhoneNumber:   request.PhoneNumber,
	}
	_, err = col.CreateDocument(ctx, doc)
	if err != nil {
		log.Printf("Error adding new user to collection: %v", err)
		return &pb.CreateAccountResponse{StatusType: 2, StatusMessage: "Error Creating User " + request.Username + ". Received the followng error " + err.Error()}, nil
	}
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "Successfully Created User " + request.Username}, nil
}
func (s *server) UpdateInfo(ctx context.Context, request *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {
	log.Printf("Updating Account Info For User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u._key"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	defer cursor.Close()
	var existingUserKey string
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUserKey)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	doc := modelUser.User{
		FirstName:     request.Username,
		MiddleInitial: request.MiddleInitial,
		LastName:      request.LastName,
		Username:      request.Username,
		Email:         request.Email,
		PhoneNumber:   request.PhoneNumber,
	}
	_, err = col.UpdateDocument(ctx, existingUserKey, doc)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return &pb.UpdateInfoResponse{StatusType: 2, StatusMessage: "Error Updating User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	return &pb.UpdateInfoResponse{StatusType: 1, StatusMessage: "Successfully Updated User " + request.Username}, nil
}
func (s *server) DeleteAccount(ctx context.Context, request *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	log.Printf("Deleting Account For User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	var query string
	if len(request.Username) > 0 {
		query = "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u._key"
	} else if len(request.Email) > 0 && len(request.Username) == 0 {
		query = "For u in users FILTER u.Email == '" + request.Email + "' RETURN u._key"
	}
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	defer cursor.Close()
	var existingUserKey string
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUserKey)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	_, err = col.RemoveDocument(dbCtx, existingUserKey)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return &pb.DeleteAccountResponse{StatusType: 2, StatusMessage: "Error Deleting User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	return &pb.DeleteAccountResponse{StatusType: 1, StatusMessage: "Successfully Deleted User " + request.Username}, nil
}
func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "This is a test message", OriginToken: "test", DerivativeToken: "test"}, nil
}
func (s *server) RefreshSession(ctx context.Context, request *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	return &pb.RefreshSessionResponse{StatusType: 1, StatusMessage: "This is a test message", DerivativeToken: "test"}, nil
}
func (s *server) PerformTwoFactorAuth(ctx context.Context, request *pb.TwoFactorAuthRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "This is a test message", OriginToken: "test", DerivativeToken: "test"}, nil
}
func (s *server) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return &pb.ChangePasswordResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	return &pb.ResetPasswordResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	return &pb.LockAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	return &pb.UnlockAccountResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) AddRole(ctx context.Context, request *pb.AddRoleRequest) (*pb.AddRoleResponse, error) {
	return &pb.AddRoleResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}
func (s *server) DeleteRole(ctx context.Context, request *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	return &pb.DeleteRoleResponse{StatusType: 1, StatusMessage: "This is a test message"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	reflection.Register(s)
	log.Printf("Started server successfully on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
