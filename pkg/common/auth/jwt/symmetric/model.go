package symmetric

import "github.com/golang-jwt/jwt/v5"

type Config struct {
	MDKey          string
	Claims         interface{ jwt.Claims }
	JwtCredentials JwtCredentials
}

type JwtCredentials struct {
	Iss        string
	Alg        string
	SigningKey string
}
