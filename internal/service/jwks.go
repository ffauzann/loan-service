package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lestrrat-go/jwx/v2/jwk"
)

func (s *service) Jwks(ctx context.Context) (jwks []*model.Jwk, err error) {
	// Iterate key pairs from config.
	for _, v := range s.config.Jwt.AsymmetricKeys {
		// Decode base64 public key.
		b, err := base64.StdEncoding.DecodeString(v.PublicKey)
		if err != nil {
			util.LogContext(ctx).Error(err.Error())
			return nil, err
		}

		// Parse public key PEM to struct.
		pub, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			util.LogContext(ctx).Error(err.Error())
			return nil, err
		}

		// Validate key length.
		modulusLength := pub.N.BitLen()
		if modulusLength < 2048 { //nolint
			util.LogContext(ctx).Error(fmt.Sprintf("one of mod has less than 2048 bits. kid: %s; len: %d", v.Kid, modulusLength))
		}

		// Construct jwks.
		jwks = append(jwks, &model.Jwk{
			KeyType:   "RSA",
			KeyID:     v.Kid,
			Usage:     "sig",
			Algorithm: "RS256",
			Modulus:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			Exponent:  base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		})
	}

	return
}
