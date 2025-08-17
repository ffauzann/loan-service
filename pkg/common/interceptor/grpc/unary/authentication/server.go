package authentication

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt/asymmetric"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/auth/jwt/symmetric"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(cfg *Config, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)

	var jwtService jwt.JwtService
	switch cfg.Alg {
	case AlgRS256:
		jwtService = jwt.NewJwtAsymmetric(&asymmetric.Config{
			MDKey:   o.mdKey,
			Claims:  o.claims,
			JwksURL: cfg.JwksURL,
			JwtCredentials: asymmetric.JwtCredentials{
				Iss: cfg.Iss,
				Alg: string(cfg.Alg),
			},
		})
	case AlgHS256:
		jwtService = jwt.NewJwtSymmetric(&symmetric.Config{
			MDKey:  o.mdKey,
			Claims: o.claims,
			JwtCredentials: symmetric.JwtCredentials{
				Iss:        cfg.Iss,
				Alg:        string(cfg.Alg),
				SigningKey: cfg.SigningKey,
			},
		})
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if slices.Contains(o.excludedMethods, util.GetMethod(info.FullMethod)) {
			return handler(ctx, req)
		}

		if jwtService == nil {
			return nil, status.Error(codes.Internal, "Internal")

		}

		ctx, err = jwtService.WithUserInfoContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Invalid or expired token")
		}

		return handler(ctx, req)
	}
}
