package symmetric

import (
	"context"
	"fmt"
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

	if tokenType, token, ok := strings.Cut(val[0], " "); ok {
		if strings.EqualFold(tokenType, expectedScheme) {
			ctx, ok = r.withUserInfoClaims(ctx, token)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "")
			}
			return ctx, nil
		}
	}

	return nil, status.Error(codes.Unauthenticated, "")
}

func (r *Config) withUserInfoClaims(ctx context.Context, token string) (context.Context, bool) {
	hmacSecret := []byte(r.JwtCredentials.SigningKey)

	jwtToken, _ := jwt.NewParser(
		jwt.WithIssuedAt(),
		jwt.WithIssuer(r.JwtCredentials.Iss),
	).ParseWithClaims(token, r.Claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.JwtCredentials.Alg {
			return nil, fmt.Errorf("Invalid signing method. expected : %v | got : %s", r.JwtCredentials.Alg, t.Header["alg"])
		}

		return hmacSecret, nil
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
