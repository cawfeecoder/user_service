package servicesUser

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	db "github.com/nfrush/user_service/user_service_server/database/mongo"
	user "github.com/nfrush/user_service/user_service_server/models/users/user"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var collection *mgo.Collection

var errorReg = regexp.MustCompile(`(index:)(?P<field>\D+)(dup key: { :)(?P<value>\D+)(})`)

func init() {
	collection = db.GetSession().C("users")
}

//CreateUser - Create and Insert a New User
func CreateUser(u user.User) error {
	log.Print(u)
	err := collection.Insert(u)
	if err != nil {
		match := errorReg.FindStringSubmatch(err.Error())
		result := make(map[string]string)
		for i, name := range errorReg.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}
		log.Printf("field: %s with value: %s", result["field"], result["value"])
		return fmt.Errorf("[ERROR] Could not create user. %s of value %s already exists!", strings.Title(strings.TrimSpace(result["field"])), strings.TrimSpace(result["value"]))
	}
	return nil
}

//FindByUsername - Find a User by their Username
func FindByUsername(username string) (user.User, error) {
	result := user.User{}
	docQuery := bson.M{"username": username}
	err := collection.Find(docQuery).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//FindByEmail - Find a User by their Email
func FindByEmail(email string) (user.User, error) {
	result := user.User{}
	docQuery := bson.M{"email": email}
	err := collection.Find(docQuery).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//FindByPhoneNumber - Find a User by their PhoneNumber
func FindByPhoneNumber(phoneNumber string) (user.User, error) {
	result := user.User{}
	docQuery := bson.M{"phoneNumber": phoneNumber}
	err := collection.Find(docQuery).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//UpdateUserInfo - Update an existing user's information
func UpdateUserInfo(u user.User) error {
	docQuery := bson.M{"username": u.Username}
	err := collection.Update(docQuery, u)
	if err != nil {
		return err
	}
	return nil
}

//DeleteUser - Deletes an existing username
func DeleteUser(u user.User) error {
	docQuery := bson.M{"username": u.Username}
	err := collection.Remove(docQuery)
	if err != nil {
		return err
	}
	return nil
}
