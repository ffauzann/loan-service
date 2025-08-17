package grpc

import (
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/service"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type srv struct {
	gen.UnimplementedAuthServiceServer
	gen.UnimplementedUserServiceServer
	gen.UnimplementedLoanServiceServer
	service service.Service
	logger  *zap.Logger
}

func New(server *grpc.Server, userSrv service.Service, logger *zap.Logger) {
	srv := srv{
		service: userSrv,
		logger:  logger,
	}
	gen.RegisterAuthServiceServer(server, &srv)
	gen.RegisterUserServiceServer(server, &srv)
	gen.RegisterLoanServiceServer(server, &srv)
	reflection.Register(server)
}
