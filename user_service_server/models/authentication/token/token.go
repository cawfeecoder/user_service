package modelToken

import "time"

//JWT Model
type JWT struct {
	Token    string `json:"token"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	IssuedAt int64  `json:"iat"`
	Expires  int64  `json:"exp:"`
	JTI      string `json:"jti"`
	Roles    string `json:"role"`
}

var signingIssuer string
var signingJTI string

//GetSigningIssuer - Returns the current issuer of tokens
func GetSigningIssuer() string {
	return signingIssuer
}

//SetSigningIssuer - Sets the current issuer of tokens
func SetSigningIssuer(issuer string) {
	signingIssuer = issuer
}

//GetSigningJTI - Returns the current JTI of tokens
func GetSigningJTI() string {
	return signingJTI
}

//SetSigningJTI - Sets the current JTI of tokens
func SetSigningJTI(jti string) {
	signingJTI = jti
}

//NewOrigin - Creates a new origin token
func NewOrigin(token string, audience string, roles string) JWT {
	newOrigin := JWT{}
	newOrigin.Token = token
	newOrigin.Issuer = GetSigningIssuer()
	newOrigin.Audience = audience
	newOrigin.IssuedAt = time.Now().Unix()
	newOrigin.Expires = time.Now().Add(time.Hour * 72).Unix()
	newOrigin.JTI = GetSigningJTI()
	newOrigin.Roles = roles
	return newOrigin
}

//NewDerivative - Creates a new derivative token
func NewDerivative(token string, audience string, roles string) JWT {
	newDeriv := JWT{}
	newDeriv.Token = token
	newDeriv.Issuer = GetSigningIssuer()
	newDeriv.Audience = audience
	newDeriv.IssuedAt = time.Now().Unix()
	newDeriv.Expires = time.Now().Add(time.Hour * 1).Unix()
	newDeriv.JTI = GetSigningJTI()
	newDeriv.Roles = roles
	return newDeriv
}

//GetToken - Returns token of JWT struct
func (j *JWT) GetToken() string {
	return j.Token
}

//SetToken - Sets token of JWT struct
func (j *JWT) SetToken(token string) {
	j.Token = token
}

//GetIssuer - Returns issuer of JWT struct
func (j *JWT) GetIssuer() string {
	return j.Token
}

//SetIssuer - Sets issuer of JWT struct
func (j *JWT) SetIssuer(issuer string) {
	j.Issuer = issuer
}

//GetAudience - Returns audience of JWT struct
func (j *JWT) GetAudience() string {
	return j.Audience
}

//SetAudience - Sets audience of JWT struct
func (j *JWT) SetAudience(aud string) {
	j.Audience = aud
}

//GetIssuedAt - Returns issued date of JWT struct
func (j *JWT) GetIssuedAt() string {
	return time.Unix(j.IssuedAt, 0).String()
}

//SetIssuedAt - Sets issued date of JWT struct
func (j *JWT) SetIssuedAt(date time.Time) {
	j.IssuedAt = date.Unix()
}

//GetExpiresAt - Returns expiry date of JWT struct
func (j *JWT) GetExpiresAt() string {
	return time.Unix(j.Expires, 0).String()
}

//SetExpiresAt - Sets expiry date of JWT struct
func (j *JWT) SetExpiresAt(date time.Time) {
	j.Expires = date.Unix()
}

//GetJTI - Returns JTI of JWT struct
func (j *JWT) GetJTI() string {
	return j.JTI
}

//SetJTI - Sets JTI of JWT struct
func (j *JWT) SetJTI(jti string) {
	j.JTI = jti
}

//GetRoles - Returns roles of JWT struct
func (j *JWT) GetRoles() string {
	return j.Roles
}

//SetRoles - Sets roles of JWT struct
func (j *JWT) SetRoles(roles string) {
	j.Roles = roles
}
