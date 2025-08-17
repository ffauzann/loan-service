package grpc

import (
	"context"
	"slices"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/constant"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/proto/gen"
)

func (s *srv) Register(ctx context.Context, req *gen.RegisterRequest) (res *gen.RegisterResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.RegisterRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Extract claims if any.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok { // Force role to be a `user` if token is not provided.
		param.User.RoleId = constant.RoleIdBorrower
	} else { // For internal SA/admin registration.
		// Validate whether claims has sufficient permission to create another admin.
		userRoleId := claims.RoleId
		if !slices.Contains(constant.AllowedRolesRegister, userRoleId) {
			err = constant.ErrPermissionDenied
			util.LogContext(ctx).Warn(err.Error())
			return
		}

		// Validate requested role_id.
		if !slices.Contains(constant.AllowedRolesRegisterMap[userRoleId], param.User.RoleId) {
			err = constant.ErrPermissionDenied
			util.LogContext(ctx).Warn(err.Error())
			return
		}
	}

	// Begin core process for the request.
	result, err := s.service.Register(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.RegisterResponse{
		Code:    gen.RegisterStatusCode(result.StatusCode),
		Reasons: result.Reasons,
	}

	return
}
