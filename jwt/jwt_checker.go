package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}
type JwtValidatorConfig struct {
	Jwks       Jwks
	Aud        string
	Issuer     string
	JwkAddress string
	cert       string
}

func NewJwtValidatorConfig(aud string, issuer string, getJwtFunc func() Jwks) *JwtValidatorConfig {
	return &JwtValidatorConfig{
		Aud:        aud,
		Issuer:     issuer,
		Jwks:       getJwtFunc(),
		JwkAddress: fmt.Sprintf("%s/.well-known/openid-configuration/jwks", issuer),
	}
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func (config *JwtValidatorConfig) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, config.ValidationKeyGetter)
}

func (config *JwtValidatorConfig) ValidationKeyGetter(token *jwt.Token) (interface{}, error) {
	audClaim, _ := token.Claims.(jwt.MapClaims)["aud"]
	validAud := false
	if real, ok := audClaim.([]interface{}); ok {
		for _, v := range real {
			if v == config.Aud {
				validAud = true
			}
		}
	} else {
		if v, ok := audClaim.(interface{}); ok {
			if v == config.Aud {
				validAud = true
			}
		}
	}
	if !validAud {
		return token, errors.New("Invalid audience.")
	}
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(config.Issuer, false)
	if !checkIss {
		return token, errors.New("Invalid issuer.")
	}

	cert, err := config.getPemCert(token)
	if err != nil {
		panic(err.Error())
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func (config *JwtValidatorConfig) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	for k, _ := range config.Jwks.Keys {
		if token.Header["kid"] == config.Jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + config.Jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}
	return cert, nil
}

func (config *JwtValidatorConfig) GetJwks() (*Jwks, error) {
	resp, err := http.Get(config.JwkAddress)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks = &Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return nil, err
	}
	return jwks, nil
}
