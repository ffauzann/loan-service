package client

import (
	"context"
	"io"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	AuthClient
	UserClient

	io.Closer
}

type UserClient interface {
	CloseAccount(ctx context.Context, req *emptypb.Empty) (res *gen.CloseAccountResponse, err error)
}

type AuthClient interface {
	Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error)
	IsUserExist(ctx context.Context, req *gen.IsUserExistRequest) (*gen.IsUserExistResponse, error)
	Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error)
	RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error)
}

type Options struct {
	GrpcAddress  string
	Interceptors []grpc.UnaryClientInterceptor
}

func New(opts Options) (Client, error) {
	conn, err := grpc.NewClient(opts.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(opts.Interceptors...),
	)
	if err != nil {
		return nil, err
	}
	return &userClient{
		conn:       conn,
		authClient: gen.NewAuthServiceClient(conn),
		userClient: gen.NewUserServiceClient(conn),
	}, nil
}

type userClient struct {
	conn       *grpc.ClientConn
	authClient gen.AuthServiceClient
	userClient gen.UserServiceClient
}

func (c *userClient) Close() error {
	return c.conn.Close()
}
