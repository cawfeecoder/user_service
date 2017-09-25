package servicesToken

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	redisDB "github.com/nfrush/user_service/user_service_server/database/redis"
	token "github.com/nfrush/user_service/user_service_server/models/authentication/token"
	user "github.com/nfrush/user_service/user_service_server/models/users/user"
)

//signingKey - Signing Key For Cookies
var signingKey = InitSigningKey()
var origClient *redis.Client
var derivClient *redis.Client

//InitSigningKey - Initalize Our Key To Sign With
func InitSigningKey() string {
	return uniuri.NewLen(32)
}

//GetSigningKey - get the current signing key
func GetSigningKey() string {
	return signingKey
}

func init() {
	token.SetSigningIssuer("Frush Development LTD")
	token.SetSigningJTI("http://example.com")
	origClient = redisDB.GetOriginClient()
	derivClient = redisDB.GetDerivClient()
}

//IssueOriginToken - Issue New JWT Token
func IssueOriginToken(u *user.User) (string, error) {
	origToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   token.GetSigningIssuer(),
		"aud":   u.GetUsername(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"jti":   token.GetSigningJTI(),
		"roles": u.GetRoles(),
	})

	origTokenString, err := origToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	_, err = origClient.Set(u.GetUsername(), origTokenString, time.Hour*72).Result()
	if err != nil {
		return "", err
	}

	fmt.Println("Issued Origin Token Successfully")
	return origTokenString, nil
}

//IssueDerivativeToken - Issue new short lived token
func IssueDerivativeToken(u *user.User) (string, error) {
	originToken, err := origClient.Get(u.GetUsername()).Result()
	if err != nil {
		return "", err
	}
	derivToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   token.GetSigningIssuer(),
		"aud":   originToken,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"jti":   token.GetSigningJTI(),
		"roles": u.GetRoles(),
	})

	derivTokenString, err := derivToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	_, err = derivClient.Set(u.GetUsername(), derivTokenString, time.Hour*1).Result()
	if err != nil {
		return "", err
	}

	fmt.Println("Issued Derivative Token Successfully")
	return derivTokenString, nil
}

//RevokeOriginToken - Revoke the Origin Token
func RevokeOriginToken(u *user.User) error {
	_, err := origClient.Del(u.GetUsername()).Result()
	if err != nil {
		return err
	}
	return nil
}

//RevokeDerivativeToken - Revoke the Derivative Token
func RevokeDerivativeToken(u *user.User) error {
	_, err := derivClient.Del(u.GetUsername()).Result()
	if err != nil {
		return err
	}
	return nil
}

//RefreshOriginToken - Reissue a new origin token
func RefreshOriginToken(u *user.User) (string, error) {
	newOrigToken, err := IssueOriginToken(u)
	if err != nil {
		return "", err
	}
	return newOrigToken, nil
}

//RefreshDerivToken - Reissue a new derivative token
func RefreshDerivToken(u *user.User) (string, error) {
	newDerivToken, err := IssueDerivativeToken(u)
	if err != nil {
		return "", err
	}
	return newDerivToken, nil
}

//OrigTokenExists - Check if Orig Token Exists
func OrigTokenExists(u *user.User) (bool, error) {
	_, err := origClient.Get(u.GetUsername()).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

//DerivTokenExists - Check if Deriv Tokens Exists
func DerivTokenExists(u *user.User) (bool, error) {
	_, err := derivClient.Get(u.GetUsername()).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

//CheckOrigToken - Check if the given originToken is the one assigned to the user
func CheckOrigToken(origToken string, u *user.User) (bool, error) {
	res, err := origClient.Get(u.GetUsername()).Result()
	if err != nil {
		return false, err
	}
	if origToken == res {
		return true, nil
	}
	return false, nil
}

//CheckDerivToken - Check if the given derivToken is the one assigned to the user
func CheckDerivToken(derivToken string, u *user.User) (bool, error) {
	res, err := derivClient.Get(u.GetUsername()).Result()
	if err != nil {
		return false, err
	}
	if derivToken == res {
		return true, nil
	}
	return false, nil
}

//GetOrigToken - Get Origin Token currently assigned to user, otherwise generate a new one
func GetOrigToken(u *user.User) (string, error) {
	res, err := origClient.Get(u.GetUsername()).Result()
	if err != nil {
		newOrigToken, err := IssueOriginToken(u)
		if err != nil {
			return "", err
		}
		return newOrigToken, nil
	}
	return res, nil
}

//GetDerivToken - Get Deriv Token currently assigned to user, otherwise generate a new one
func GetDerivToken(u *user.User) (string, error) {
	res, err := derivClient.Get(u.GetUsername()).Result()
	if err != nil {
		newDerivToken, err := IssueDerivativeToken(u)
		if err != nil {
			return "", err
		}
		return newDerivToken, nil
	}
	return res, nil
}
