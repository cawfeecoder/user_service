package modelUser

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//Role model
type Role struct {
	ID   int
	Name string
}

//User model
type User struct {
	FirstName     string
	MiddleInitial string
	LastName      string
	Username      string
	Password      string
	Status        int
	TimeLocked    time.Time
	Email         string
	PhoneNumber   string
	Roles         []Role
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
	newUser.Status = 2
	newUser.Roles = []Role{Role{ID: 1, Name: "User"}}
	return newUser
}

func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
