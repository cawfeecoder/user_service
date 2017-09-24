package main

import (
	"log"
	"net"
	"time"

	driver "github.com/arangodb/go-driver"
	pb "github.com/nfrush/user_service/user_service_pb"
	db "github.com/nfrush/user_service/user_service_server/database/arangodb"
	redis "github.com/nfrush/user_service/user_service_server/database/redis"
	"github.com/nfrush/user_service/user_service_server/models/user"
	"github.com/nfrush/user_service/user_service_server/services/token"
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
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	doc := modelUser.New(request.FirstName, request.MiddleInitial, request.LastName, request.Username, request.Password, request.Email, request.PhoneNumber)
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
	log.Printf("Logging In Account For User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	defer cursor.Close()
	var existingUser modelUser.User
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	autherr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(request.Password))
	if autherr != nil {
		log.Printf("Error authenticating user: %v", err)
		return &pb.LoginResponse{StatusType: 2, StatusMessage: "Error Authenticating User " + request.Username}, nil
	}
	client := redis.GetClient()
	val, err := client.Get(request.Username).Result()
	log.Printf("Read value: %s", val)
	if err != nil {
		originToken, _ := servicesToken.IssueOriginToken(request.Username, existingUser.Roles)
		derivativeToken, _ := servicesToken.IssueDerivativeToken(originToken, existingUser.Roles)
		client.Set(request.Username, originToken, 72*time.Hour)
		client.Set(originToken, derivativeToken, time.Hour)
		return &pb.LoginResponse{StatusType: 1, StatusMessage: "Successfully Authenticated User " + request.Username, OriginToken: originToken, DerivativeToken: derivativeToken}, nil
	}
	derivativeToken, _ := servicesToken.IssueDerivativeToken(val, existingUser.Roles)
	client.Set(val, derivativeToken, time.Hour)
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "Successfully Authenticated User " + request.Username, OriginToken: val, DerivativeToken: derivativeToken}, nil
}
func (s *server) RefreshSession(ctx context.Context, request *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	client := redis.GetClient()
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	defer cursor.Close()
	var existingUser modelUser.User
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	val, err := client.Get(request.Username).Result()
	log.Printf("Read value: %s", val)
	if err != nil {
		originToken, _ := servicesToken.IssueOriginToken(existingUser.Username, existingUser.Roles)
		derivativeToken, _ := servicesToken.IssueDerivativeToken(originToken, existingUser.Roles)
		return &pb.RefreshSessionResponse{StatusType: 1, StatusMessage: "Origin token and Derivative Token was succesfully refreshed. User has been issued a new session", OriginToken: originToken, DerivativeToken: derivativeToken}, nil
	}
	if val != request.OriginToken {
		log.Printf("Origin token mismatch. Origin token provided in request does not match the issued token")
		return &pb.RefreshSessionResponse{StatusType: 2, StatusMessage: "Could not refresh token. Origin token mismatch"}, nil
	}
	derivativeToken, _ := servicesToken.IssueDerivativeToken(request.OriginToken, existingUser.Roles)
	client.Set(val, derivativeToken, time.Hour)
	return &pb.RefreshSessionResponse{StatusType: 1, StatusMessage: "Successfully Refreshed Session For User " + request.Username, DerivativeToken: derivativeToken}, nil
}
func (s *server) PerformTwoFactorAuth(ctx context.Context, request *pb.TwoFactorAuthRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "Two Factor Authenication is currently Disabled", OriginToken: "", DerivativeToken: ""}, nil
}
func (s *server) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Printf("Logging Out Account For User %s", request.Username)
	client := redis.GetClient()
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	defer cursor.Close()
	var existingUser modelUser.User
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	origin, err := client.Get(request.Username).Result()
	log.Printf("Read value: %s", origin)
	if err != nil {
		log.Printf("User session has already expired")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "User session has already expired. Cannot log out"}, nil
	}
	if origin != request.OriginToken {
		log.Printf("Origin Token Mismatch")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "Origin Token Mismatch"}, nil
	}
	deriv, err := client.Get(request.OriginToken).Result()
	log.Printf("Read value: %s", origin)
	if err != nil {
		log.Printf("User session has already expired")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "User session has already expired. Cannot log out"}, nil
	}
	if deriv != request.DerivativeToken {
		log.Printf("Deriative Token Mismatch")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "Derivative Token Mismatch"}, nil
	}
	return &pb.LogoutResponse{StatusType: 1, StatusMessage: "Successfully Logged Out User " + request.Username}, nil
}
func (s *server) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	log.Printf("Changing Password For User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Printf("Error getting collection.")
	}
	defer cursor.Close()
	var existingUser driver.DocumentMeta
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	var existingUserDocument modelUser.User
	_, err = col.ReadDocument(dbCtx, existingUser.Key, &existingUserDocument)
	if err != nil {
		log.Printf("Error getting user document")
	}
	autherr := bcrypt.CompareHashAndPassword([]byte(existingUserDocument.Password), []byte(request.OldPassword))
	if autherr != nil {
		log.Printf("Error authenticating user: %v", err)
		return &pb.ChangePasswordResponse{StatusType: 2, StatusMessage: "Error Authenticating User " + request.Username}, nil
	}
	existingUserDocument.Password = request.NewPassword
	_, err = col.UpdateDocument(dbCtx, existingUser.Key, existingUserDocument)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.ChangePasswordResponse{StatusType: 2, StatusMessage: "Error Changing Password."}, nil
	}
	return &pb.ChangePasswordResponse{StatusType: 1, StatusMessage: "Successfully Changed Password For User " + request.Username}, nil
}
func (s *server) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	log.Printf("Putting Status to PASSWORD_RESET for User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Printf("Error getting collection.")
	}
	defer cursor.Close()
	var existingUser driver.DocumentMeta
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	var existingUserDocument modelUser.User
	_, err = col.ReadDocument(dbCtx, existingUser.Key, &existingUserDocument)
	if err != nil {
		log.Printf("Error getting user document")
	}
	existingUserDocument.Status = 5
	_, err = col.UpdateDocument(dbCtx, existingUser.Key, existingUserDocument)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.ResetPasswordResponse{StatusType: 2, StatusMessage: "Error Putting User in Reset Password Status."}, nil
	}
	return &pb.ResetPasswordResponse{StatusType: 1, StatusMessage: "Successfully Set PASSWORD_RESET Status For User " + request.Username}, nil
}
func (s *server) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	log.Printf("Putting Status to LOCK for User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Printf("Error getting collection.")
	}
	defer cursor.Close()
	var existingUser driver.DocumentMeta
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	var existingUserDocument modelUser.User
	_, err = col.ReadDocument(dbCtx, existingUser.Key, &existingUserDocument)
	if err != nil {
		log.Printf("Error getting user document: %v", err)
	}
	if request.LockType == 1 {
		existingUserDocument.Status = 3
		existingUserDocument.TimeLocked = time.Now().Add(time.Minute * 15)
	} else if request.LockType == 2 {
		existingUserDocument.Status = 4
		existingUserDocument.TimeLocked = time.Date(9999, 12, 30, 23, 59, 59, 0, time.UTC)
	}
	_, err = col.UpdateDocument(dbCtx, existingUser.Key, existingUserDocument)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.LockAccountResponse{StatusType: 2, StatusMessage: "Error Putting User in Lock Status."}, nil
	}
	return &pb.LockAccountResponse{StatusType: 1, StatusMessage: "Successfully Set Lock Status For User " + request.Username}, nil
}
func (s *server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	log.Printf("Removing LOCK status for User %s", request.Username)
	dbCtx := context.Background()
	db, err := db.GetClient().Database(dbCtx, "test")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err != nil {
		log.Fatalf("Could not find collection: %v", err)
	}
	query := "FOR u in users FILTER u.Username == '" + request.Username + "' RETURN u"
	cursor, err := db.Query(dbCtx, query, nil)
	if err != nil {
		log.Fatalf("Encountered this error querying: %v", err)
	}
	col, err := db.Collection(dbCtx, "users")
	if err != nil {
		log.Printf("Error getting collection.")
	}
	defer cursor.Close()
	var existingUser driver.DocumentMeta
	for {
		_, cerr := cursor.ReadDocument(dbCtx, &existingUser)
		if driver.IsNoMoreDocuments(cerr) {
			break
		} else if cerr != nil {
			log.Fatalf("Encounted this error while reading cursor: %v", err)
		}
	}
	var existingUserDocument modelUser.User
	_, err = col.ReadDocument(dbCtx, existingUser.Key, &existingUserDocument)
	if err != nil {
		log.Printf("Error getting user document: %v", err)
	}
	existingUserDocument.Status = 1
	existingUserDocument.TimeLocked = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err = col.UpdateDocument(dbCtx, existingUser.Key, existingUserDocument)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.UnlockAccountResponse{StatusType: 2, StatusMessage: "Error Unlocking User."}, nil
	}
	return &pb.UnlockAccountResponse{StatusType: 1, StatusMessage: "Successfully Unlocked User " + request.Username}, nil
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
