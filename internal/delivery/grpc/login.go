package grpc

import (
	"context"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"
)

func (s *srv) Login(ctx context.Context, req *gen.LoginRequest) (res *gen.LoginResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.LoginRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Begin core process for the request.
	result, err := s.service.Login(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.LoginResponse{
		AccessToken:  result.Token.AccessToken,
		RefreshToken: result.Token.RefreshToken,
	}

	return
}
