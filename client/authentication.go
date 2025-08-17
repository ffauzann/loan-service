package client

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"
)

func (c *userClient) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	return c.authClient.Register(ctx, req)
}

func (c *userClient) IsUserExist(ctx context.Context, req *gen.IsUserExistRequest) (*gen.IsUserExistResponse, error) {
	return c.authClient.IsUserExist(ctx, req)
}

func (c *userClient) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	return c.authClient.Login(ctx, req)
}

func (c *userClient) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
	return c.authClient.RefreshToken(ctx, req)
}
