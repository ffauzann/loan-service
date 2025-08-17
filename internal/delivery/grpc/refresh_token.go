package grpc

import (
	"context"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/ffauzann/loan-service/proto/gen"
)

func (s *srv) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (res *gen.RefreshTokenResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.RefreshTokenRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Begin core process for the request.
	result, err := s.service.RefreshToken(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = util.CastStruct[gen.RefreshTokenResponse](result)
	return
}
