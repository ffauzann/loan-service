package grpc

import (
	"github.com/ffauzann/loan-service/internal/service"
	"github.com/ffauzann/loan-service/proto/gen"

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
