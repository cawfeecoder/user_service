package servicesToken

import (
	"errors"
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/dgrijalva/jwt-go"
	"github.com/nfrush/user_service/user_service_server/models/user"
)

//signingKey - Signing Key For Cookies
var signingKey = InitSigningKey()

//InitSigningKey - Initalize Our Key To Sign With
func InitSigningKey() string {
	return uniuri.NewLen(32)
}

//GetSigningKey - get the current signing key
func GetSigningKey() string {
	return signingKey
}

//IssueOriginToken - Issue New JWT Token
func IssueOriginToken(username string, roles []modelUser.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   "Frush Development LTD",
		"aud":   username,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"jti":   "http://example.com",
		"roles": roles,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	//issuedToken := modelToken.JWT{Token: tokenString, Issuer: "Frush Development LTD", Audience: u.Username, IssuedAt: time.Now().Unix(), Expires: time.Now().Add(time.Hour * 72).Unix(), JTI: "http://example.com"}

	fmt.Println("Issued Token Successfully")
	return tokenString, nil
}

//IssueDerivativeToken - Issue new short lived token
func IssueDerivativeToken(originToken string, roles []modelUser.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   "Frush Development LTD",
		"aud":   originToken,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"jti":   "http://example.com",
		"roles": roles,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	//issuedToken := modelToken.JWT{Token: tokenString, Issuer: "Frush Development LTD", Audience: originToken, IssuedAt: time.Now().Unix(), Expires: time.Now().Add(time.Hour * 1).Unix(), JTI: "http://example.com"}

	fmt.Println("Issued Token Successfully")
	return tokenString, nil
}

//RevokeToken - Revoke the JWT Token
func RevokeToken(u *modelUser.User) error {
	return nil
}

//RefreshToken - Reissue a new token
func RefreshToken(u *modelUser.User) (string, error) {
	return "", nil
}

//TokenExists - Check if Token Exists
func TokenExists(token string) (bool, error) {
	return false, nil
}

//TokenExistsUser - Checks if a user has an assigned token
func TokenExistsUser(u *modelUser.User) (bool, error) {
	return false, nil
}

//RequiresAuth - Authenicates user on service
func RequiresAuth(token string) (bool, error) {
	return false, errors.New("The Token Does Not Exist")
}
