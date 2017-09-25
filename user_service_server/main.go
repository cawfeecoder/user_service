package main

import (
	"log"
	"net"
	"time"

	pb "github.com/nfrush/user_service/user_service_pb"
	user "github.com/nfrush/user_service/user_service_server/models/users/user"
	"github.com/nfrush/user_service/user_service_server/services/token"
	"github.com/nfrush/user_service/user_service_server/services/users"
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
	newUser := user.New(request.FirstName, request.MiddleInitial, request.LastName, request.Username, request.Password, request.Email, request.PhoneNumber)
	err := servicesUser.CreateUser(newUser)
	if err != nil {
		log.Printf("Error encountered creating new user: %s", err.Error())
		return &pb.CreateAccountResponse{StatusType: 2, StatusMessage: "Error Creating User " + request.Username + ". Received the followng error: " + err.Error()}, nil
	}
	return &pb.CreateAccountResponse{StatusType: 1, StatusMessage: "Successfully Created User " + request.Username}, nil
}
func (s *server) UpdateInfo(ctx context.Context, request *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {
	log.Printf("Updating Account Info For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Could not find user: %v", err)
		return &pb.UpdateInfoResponse{StatusType: 3, StatusMessage: "Error Update User " + request.Username + ". Could not find user."}, nil
	}
	if len(request.FirstName) > 0 {
		existingUser.FirstName = request.FirstName
	}
	if len(request.MiddleInitial) > 0 {
		existingUser.MiddleInitial = request.MiddleInitial
	}
	if len(request.LastName) > 0 {
		existingUser.LastName = request.LastName
	}
	if len(request.Email) > 0 {
		existingUser.Email = request.Email
	}
	if len(request.PhoneNumber) > 0 {
		existingUser.PhoneNumber = request.PhoneNumber
	}
	err = servicesUser.UpdateUserInfo(existingUser)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return &pb.UpdateInfoResponse{StatusType: 2, StatusMessage: "Error Updating User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	return &pb.UpdateInfoResponse{StatusType: 1, StatusMessage: "Successfully Updated User " + request.Username}, nil
}
func (s *server) DeleteAccount(ctx context.Context, request *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	log.Printf("Deleting Account For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Could not find user: %v", err)
		return &pb.DeleteAccountResponse{StatusType: 3, StatusMessage: "Error Deleting User " + request.Username + ". Could not find user."}, nil
	}
	err = servicesUser.DeleteUser(existingUser)
	if err != nil {
		log.Printf("Could not delete user: %v", err)
		return &pb.DeleteAccountResponse{StatusType: 2, StatusMessage: "Error Deleting User " + request.Username + ". Received the following error " + err.Error()}, nil
	}
	return &pb.DeleteAccountResponse{StatusType: 1, StatusMessage: "Successfully Deleted User " + request.Username}, nil
}
func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Logging In Account For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.LoginResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	if existingUser.GetStatus() == "ACCOUNT_LOCK" {
		return &pb.LoginResponse{StatusType: 9, StatusMessage: "Account is permenantly locked. Please contact your local administrator"}, nil
	}
	if existingUser.LockedUntil.After(time.Now()) {
		log.Printf("User account is still locked for: " + time.Now().Sub(existingUser.LockedUntil).String())
		return &pb.LoginResponse{StatusType: 4, StatusMessage: "User account is stilled locked for " + existingUser.LockedUntil.Sub(time.Now()).Round(time.Second).String()}, nil
	}
	auth := existingUser.CheckPassword(request.Password)
	if auth == false {
		log.Printf("Error authenticating user: %v", err)
		if existingUser.FailedLogins >= 8 {
			existingUser.LockAccount(2)
			err = servicesUser.UpdateUserInfo(existingUser)
			if err != nil {
				log.Printf("Error locking user account")
				return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
			}
			return &pb.LoginResponse{StatusType: 5, StatusMessage: "Too many failed login attempts. User account is now locked. Please contact your local administrator"}, nil
		} else if existingUser.FailedLogins < 8 && existingUser.FailedLogins%3 == 2 {
			existingUser.LockAccount(1)
			err = servicesUser.UpdateUserInfo(existingUser)
			if err != nil {
				log.Printf("Error locking user account")
				return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
			}
			return &pb.LoginResponse{StatusType: 7, StatusMessage: "Too many failed login attempts. Please wait 15 minutes before attempting to login again."}, nil
		} else {
			existingUser.LockAccount(1)
			err = servicesUser.UpdateUserInfo(existingUser)
			if err != nil {
				log.Printf("Error locking user account")
				return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
			}
			return &pb.LoginResponse{StatusType: 2, StatusMessage: "Error Authenticating User " + request.Username + "."}, nil
		}
	}
	existingUser.UnlockAccount()
	err = servicesUser.UpdateUserInfo(existingUser)
	if err != nil {
		log.Printf("Error resetting failed login attempts for user account")
		return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
	}
	originToken, err := servicesToken.GetOrigToken(&existingUser)
	if err != nil {
		log.Printf("Error getting origin token or issuing a new origin token.")
		return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
	}
	derivToken, err := servicesToken.GetDerivToken(&existingUser)
	if err != nil {
		log.Printf("Error getting derivative token or issuing a new derivative token.")
		return &pb.LoginResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
	}
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "Successfully Authenticated User " + request.Username, OriginToken: originToken, DerivativeToken: derivToken}, nil
}
func (s *server) RefreshSession(ctx context.Context, request *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	log.Printf("Refreshing Session For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.RefreshSessionResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	originToken, err := servicesToken.GetOrigToken(&existingUser)
	if err != nil {
		log.Printf("Error getting origin token or issuing a new origin token.")
		return &pb.RefreshSessionResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
	}
	derivToken, err := servicesToken.GetDerivToken(&existingUser)
	if err != nil {
		log.Printf("Error getting derivative token or issuing a new derivative token.")
		return &pb.RefreshSessionResponse{StatusType: 6, StatusMessage: "Unrecoverable Error Occured. Could not process login. Please try again or contact your local administrator"}, nil
	}
	if originToken != request.OriginToken {
		log.Printf("Origin token mismatch. Origin token provided in request does not match the issued token")
		return &pb.RefreshSessionResponse{StatusType: 2, StatusMessage: "Could not refresh token. Origin token mismatch"}, nil
	}
	return &pb.RefreshSessionResponse{StatusType: 1, StatusMessage: "Successfully Refreshed Session For User " + request.Username, OriginToken: originToken, DerivativeToken: derivToken}, nil
}
func (s *server) PerformTwoFactorAuth(ctx context.Context, request *pb.TwoFactorAuthRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{StatusType: 1, StatusMessage: "Two Factor Authenication is currently Disabled", OriginToken: "", DerivativeToken: ""}, nil
}
func (s *server) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Printf("Logging Out Account For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.LogoutResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	origin, err := servicesToken.OrigTokenExists(&existingUser)
	if origin == false || err != nil {
		log.Printf("User session has already expired")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "User session has already expired. Cannot log out"}, nil
	}
	originToken, err := servicesToken.CheckOrigToken(request.OriginToken, &existingUser)
	if err != nil {
		log.Printf("Could not check validity of origin token")
		return &pb.LogoutResponse{StatusType: 4, StatusMessage: "Unrecoverable Error Occured. Could not check validity of origin token. Please try again or contact your local administrator"}, nil
	}
	if originToken == false {
		log.Printf("Origin Token Mismatch")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "Origin Token Mismatch"}, nil
	}
	deriv, err := servicesToken.DerivTokenExists(&existingUser)
	if deriv == false || err != nil {
		log.Printf("User session has already expired")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "User session has already expired. Cannot log out"}, nil
	}
	derivToken, err := servicesToken.CheckDerivToken(request.DerivativeToken, &existingUser)
	if err != nil {
		log.Printf("Could not check validity of derivative token")
		return &pb.LogoutResponse{StatusType: 4, StatusMessage: "Unrecoverable Error Occured. Could not check validity of derivative token. Please try again or contact your local administrator"}, nil
	}
	if derivToken == false {
		log.Printf("Deriative Token Mismatch")
		return &pb.LogoutResponse{StatusType: 2, StatusMessage: "Derivative Token Mismatch"}, nil
	}
	err = servicesToken.RevokeDerivativeToken(&existingUser)
	if err != nil {
		log.Printf("Failed to delete Derivative Token in Redis due to: %s", err.Error())
		return &pb.LogoutResponse{StatusType: 4, StatusMessage: "User session failed to expire. Derivative Token could not be invalidated"}, nil
	}
	err = servicesToken.RevokeOriginToken(&existingUser)
	if err != nil {
		log.Printf("Failed to delete Origin Token in Redis due to: %s", err.Error())
		return &pb.LogoutResponse{StatusType: 4, StatusMessage: "User session failed to expire. Origin Token could not be invalidated"}, nil
	}
	return &pb.LogoutResponse{StatusType: 1, StatusMessage: "Successfully Logged Out User " + request.Username}, nil
}
func (s *server) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	log.Printf("Changing Password For User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.ChangePasswordResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	auth := existingUser.CheckPassword(request.OldPassword)
	if auth == false {
		log.Printf("Error authenticating user: %v", err)
		return &pb.ChangePasswordResponse{StatusType: 2, StatusMessage: "Error Authenticating User " + request.Username}, nil
	}
	existingUser.SetPassword(request.NewPassword)
	err = servicesUser.UpdateUserInfo(existingUser)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.ChangePasswordResponse{StatusType: 2, StatusMessage: "Error Changing Password."}, nil
	}
	return &pb.ChangePasswordResponse{StatusType: 1, StatusMessage: "Successfully Changed Password For User " + request.Username}, nil
}
func (s *server) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	log.Printf("Putting Status to PASSWORD_RESET for User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.ResetPasswordResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	existingUser.ResetPassword()
	err = servicesUser.UpdateUserInfo(existingUser)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.ResetPasswordResponse{StatusType: 2, StatusMessage: "Error Putting User in Reset Password Status."}, nil
	}
	return &pb.ResetPasswordResponse{StatusType: 1, StatusMessage: "Successfully Set PASSWORD_RESET Status For User " + request.Username}, nil
}
func (s *server) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	log.Printf("Putting Status to LOCK for User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.LockAccountResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	existingUser.LockAccount(int(request.LockType))
	err = servicesUser.UpdateUserInfo(existingUser)
	if err != nil {
		log.Printf("Receive error trying to update password: %v", err)
		return &pb.LockAccountResponse{StatusType: 2, StatusMessage: "Error Putting User in Lock Status."}, nil
	}
	return &pb.LockAccountResponse{StatusType: 1, StatusMessage: "Successfully Set Lock Status For User " + request.Username}, nil
}
func (s *server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	log.Printf("Removing LOCK status for User %s", request.Username)
	existingUser, err := servicesUser.FindByUsername(request.Username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return &pb.UnlockAccountResponse{StatusType: 3, StatusMessage: "Error Finding User " + request.Username + ". Received the following error: " + err.Error()}, nil
	}
	existingUser.UnlockAccount()
	err = servicesUser.UpdateUserInfo(existingUser)
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
