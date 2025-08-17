package grpc

import (
	"context"
	"regexp"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/constant"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util/str"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *srv) IsUserExist(ctx context.Context, req *gen.IsUserExistRequest) (res *gen.IsUserExistResponse, err error) {
	if err = validateIsUserExist(req); err != nil {
		return
	}

	param := util.CastStruct[model.IsUserExistRequest](req)
	param.PhoneNumber = str.PhoneWithCountryCode(param.PhoneNumber, constant.DefaultCountryCode, true)
	result, err := s.service.IsUserExist(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	res = util.CastStruct[gen.IsUserExistResponse](result)
	return
}

func validateIsUserExist(req *gen.IsUserExistRequest) error {
	if req.Email != "" {
		regexEmail := regexp.MustCompile(constant.RegexEmail)
		if !regexEmail.MatchString(req.GetEmail()) {
			return constant.ErrMalformedEmail
		}
	}

	if len(req.GetPhoneNumber()) < 4 { //nolint
		req.PhoneNumber = ""
	}

	if req.Email == "" && req.PhoneNumber == "" {
		return constant.ErrNoArg
	}

	return nil
}

func (s *srv) CloseAccount(ctx context.Context, req *emptypb.Empty) (res *gen.CloseAccountResponse, err error) {
	// Cast and validate request.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok {
		err = constant.ErrUnauthenticated
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Begin core process for the request.
	result, err := s.service.CloseAccount(ctx, &model.CloseAccountRequest{
		UserId: claims.UserId,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = util.CastStruct[gen.CloseAccountResponse](result)

	return
}
