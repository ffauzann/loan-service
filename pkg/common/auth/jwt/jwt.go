package jwt

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt/asymmetric"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt/symmetric"
)

type UserInfoKey struct{}

type JwtService interface {
	// WithUserInfoContext returns new context with user info key.
	WithUserInfoContext(ctx context.Context) (context.Context, error)
	// UserInfoFromContext returns claims from given context.
	// UserInfoFromContext(ctx context.Context) (interface{ jwt.Claims }, bool)
}

func NewJwtAsymmetric(cfg *asymmetric.Config) JwtService {
	return cfg
}

func NewJwtSymmetric(cfg *symmetric.Config) JwtService {
	return cfg
}
