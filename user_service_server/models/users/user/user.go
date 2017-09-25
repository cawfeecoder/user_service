package modelUser

import (
	"log"
	"time"

	role "github.com/nfrush/user_service/user_service_server/models/users/role"
	status "github.com/nfrush/user_service/user_service_server/models/users/status"
	"golang.org/x/crypto/bcrypt"
)

//User model
type User struct {
	FirstName       string        `bson:"firstName"`
	MiddleInitial   string        `bson:"middleInitial"`
	LastName        string        `bson:"lastName"`
	Username        string        `bson:"username"`
	Password        string        `bson:"password"`
	Status          status.Status `bson:"status"`
	LastLogin       time.Time     `bson:"lastLogin"`
	LastFailedLogin time.Time     `bson:"lastFailedLogin"`
	FailedLogins    int           `bson:"failedLogins"`
	PasswordExpire  time.Time     `bson:"passwordExpire"`
	LockedUntil     time.Time     `bson:"lockedUntil"`
	Email           string        `bson:"email"`
	PhoneNumber     string        `bson:"phoneNumber"`
	Roles           []role.Role   `bson:"roles"`
	LastUpdatedBy   string        `bson:"lastUpdatedBy"`
	LastUpdatedOn   time.Time     `bson:"lastUpdatedOn"`
}

//New - Constructs a new user from a full set of user information
func New(firstName string, middleInitial string, lastName string, username string, password string, email string, phonenumber string) User {
	newUser := User{}
	hashedPassword, err := generateHash(password)
	if err != nil {
		log.Panicf("Could not generate hashed password: %v", err)
	}
	newUser.FirstName = firstName
	newUser.MiddleInitial = middleInitial
	newUser.LastName = lastName
	newUser.Username = username
	newUser.Password = hashedPassword
	newUser.Email = email
	newUser.PhoneNumber = phonenumber
	newUser.Status = status.UNVERIFIED
	newUser.Roles = []role.Role{role.USER}
	newUser.LastUpdatedBy = "System"
	newUser.LastUpdatedOn = time.Now()
	return newUser
}

func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//GetFirstName - Returns firstName of User struct
func (u *User) GetFirstName() string {
	return u.FirstName
}

//SetFirstName - Sets firstName of User struct
func (u *User) SetFirstName(firstName string) {
	u.FirstName = firstName
}

//GetMiddleInitial - Returns middleInitial of User struct
func (u *User) GetMiddleInitial() string {
	return u.MiddleInitial
}

//SetMiddleInitial Returns middleInitial of User struct
func (u *User) SetMiddleInitial(middleInitial string) {
	u.MiddleInitial = middleInitial
}

//GetLastName - Returns lastName of User struct
func (u *User) GetLastName() string {
	return u.LastName
}

//SetLastName - Sets lastName of User struct
func (u *User) SetLastName(lastName string) {
	u.LastName = lastName
}

//GetFullName - Returns a formatted full name of User struct
func (u *User) GetFullName(format string) string {
	if format == "Military" {
		return u.LastName + ", " + u.FirstName + " " + u.MiddleInitial
	}
	return u.FirstName + " " + u.MiddleInitial + " " + u.LastName
}

//GetUsername - Returns username of User struct
func (u *User) GetUsername() string {
	return u.Username
}

//SetUsername - Sets username of User struct
func (u *User) SetUsername(username string) {
	u.Username = username
}

// SetPassword - Sets password of User struct
func (u *User) SetPassword(password string) {
	hashedPassword, err := generateHash(password)
	if err != nil {
		log.Panicf("Could not generate hashed password: %v", err)
	}
	u.Password = hashedPassword
}

// CheckPassword - Checks the user's password against a presented password for authentication
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

//GetEmail - Returns email of User struct
func (u *User) GetEmail() string {
	return u.Email
}

//SetEmail - Sets email of User struct
func (u *User) SetEmail(email string) {
	u.Email = email
}

//GetPhoneNumber - Returns phone number of User struct
func (u *User) GetPhoneNumber() string {
	return u.PhoneNumber
}

//SetPhoneNumber - Sets phone number of User struct
func (u *User) SetPhoneNumber(phoneNumber string) {
	u.PhoneNumber = phoneNumber
}

// GetStatus - Returns status of User struct
func (u *User) GetStatus() string {
	return u.Status.GetName()
}

// LockAccount - Locks user based on type {1 = FAILED LOGIN, 2 = ACCOUNT LOCK}
func (u *User) LockAccount(lock int) {
	switch lock {
	case 1:
		if u.FailedLogins%3 == 2 {
			u.Status = status.PASSWORD_LOCK
			u.LockedUntil = time.Now().Add(time.Minute * 15)
		}
		u.FailedLogins++
	case 2:
		u.Status = status.ACCOUNT_LOCK
		u.LockedUntil = time.Date(9999, 12, 30, 23, 59, 59, 0, time.UTC)
		u.FailedLogins = 0
	default:
		u.Status = status.ACCOUNT_LOCK
		u.LockedUntil = time.Date(9999, 12, 30, 23, 59, 59, 0, time.UTC)
		u.FailedLogins = 0
	}
}

//UnlockAccount - Unlocks user
func (u *User) UnlockAccount() {
	u.Status = status.ACTIVE
	u.LockedUntil = time.Time{}
	u.FailedLogins = 0
}

//Verify - Verify user email or phone number and active account
func (u *User) Verify() {
	u.Status = status.ACTIVE
}

//ResetPassword - Set user status to PASSWORD_RESET
func (u *User) ResetPassword() {
	u.Status = status.PASSWORD_RESET
}

// GetRoles - Returns roles of User struct
func (u *User) GetRoles() []string {
	var roleNames []string
	for _, v := range u.Roles {
		roleNames = append(roleNames, v.GetName())
	}
	return roleNames
}

// AddRole - Adds a single role to User struct
func (u *User) AddRole(role role.Role) {
	u.Roles = append(u.Roles, role)
}

// AddRoles - Adds an array of roles to User struct
func (u *User) AddRoles(roles []role.Role) {
	for _, v := range roles {
		u.Roles = append(u.Roles, v)
	}
}

// RemoveRole - Removes a single role of User struct
func (u *User) RemoveRole(role role.Role) {
	for i, v := range u.Roles {
		if v.GetID() == role.GetID() {
			u.Roles[i] = u.Roles[len(u.Roles)-1]
		}
	}
	u.Roles = u.Roles[:len(u.Roles)-1]
}

// Deprecate - Removes all roles of User struct
func (u *User) Deprecate() {
	u.Roles = []role.Role{role.USER}
}

// GetLastUpdated - Gets lastUpdated fields of User struct
func (u *User) GetLastUpdated() string {
	return u.LastUpdatedBy + " on " + u.LastUpdatedOn.String()
}

// SetLastUpdated - Sets lastUpdated fields of User struct
func (u *User) SetLastUpdated(user string, date time.Time) {
	u.LastUpdatedBy = user
	u.LastUpdatedOn = date
}
