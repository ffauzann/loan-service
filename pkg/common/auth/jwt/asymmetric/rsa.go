package asymmetric

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt/ctxval"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const expectedScheme = "Bearer"

func (r *Config) WithUserInfoContext(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	val := md.Get(r.MDKey)
	if len(val) == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	if tokenType, tokenString, ok := strings.Cut(val[0], " "); ok {
		if strings.EqualFold(tokenType, expectedScheme) {
			ctx, ok = r.withUserInfoClaims(ctx, tokenString)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "")
			}
			return ctx, nil
		}
	}

	return nil, status.Error(codes.Unauthenticated, "")
}

func (r *Config) withUserInfoClaims(ctx context.Context, tokenString string) (context.Context, bool) {
	if err := r.getJwks(ctx, r); err != nil {
		return ctx, false
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, r.Claims)
	if err != nil {
		return ctx, false
	}

	var publicKey *rsa.PublicKey
	for _, v := range r.jwks {
		if v.KeyID == token.Header["kid"] {
			nBytes, err := base64.RawURLEncoding.DecodeString(v.Modulus)
			if err != nil {
				return ctx, false
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(v.Exponent)
			if err != nil {
				return ctx, false
			}

			publicKey = &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: int(new(big.Int).SetBytes(eBytes).Int64()),
			}
			break
		}
	}

	jwtToken, _ := jwt.NewParser(
		jwt.WithIssuedAt(),
		jwt.WithIssuer(r.JwtCredentials.Iss),
	).ParseWithClaims(tokenString, r.Claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.JwtCredentials.Alg {
			return nil, fmt.Errorf("Invalid signing method. expected : %v | got : %s", r.JwtCredentials.Alg, t.Header["alg"])
		}

		return publicKey, nil
	})

	if jwtToken == nil {
		return ctx, false
	}

	val, ok := jwtToken.Claims.(interface{ jwt.Claims })
	if ok && jwtToken.Valid {
		return ctxval.SetUserInfo(ctx, val), true
	}

	return ctx, false
}
