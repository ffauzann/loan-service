package asymmetric

import "github.com/golang-jwt/jwt/v5"

type Config struct {
	MDKey          string
	Claims         interface{ jwt.Claims }
	JwtCredentials JwtCredentials
	JwksURL        string

	jwks []*Jwk
}

type JwtCredentials struct {
	Iss        string
	Alg        string
	SigningKey string
}

type Jwk struct {
	KeyType   string `json:"kty"`
	KeyID     string `json:"kid"`
	Usage     string `json:"use"`
	Algorithm string `json:"alg"`
	Modulus   string `json:"n"`
	Exponent  string `json:"e"`
}
