package grpc

import (
	"context"
	"slices"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/ffauzann/loan-service/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateLoan handles the creation of a new loan proposal.
// nolint
func (s *srv) CreateLoan(ctx context.Context, req *gen.CreateLoanRequest) (res *gen.CreateLoanResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.CreateLoanRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Extract claims from context.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	if !slices.Contains(constant.AllowedRolesProposeLoan, claims.RoleId) {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Set BorrowerId from claims.
	param.BorrowerId = claims.UserId

	// Begin core process for the request.
	result, err := s.service.CreateLoan(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.CreateLoanResponse{
		LoanId:    result.LoanId,
		State:     result.State,
		CreatedAt: timestamppb.New(result.CreatedAt),
	}

	return
}

// ApproveLoan handles the approval of a loan proposal.
// nolint
func (s *srv) ApproveLoan(ctx context.Context, req *gen.ApproveLoanRequest) (res *gen.ApproveLoanResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.ApproveLoanRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Extract claims from context.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate role_id from claims.
	if !slices.Contains(constant.AllowedRolesApproveLoan, claims.RoleId) {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Set ValidatorId from claims.
	param.ValidatorId = claims.UserId

	// Begin core process for the request.
	result, err := s.service.ApproveLoan(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.ApproveLoanResponse{
		LoanId:       result.LoanId,
		State:        result.State,
		ApprovalDate: timestamppb.New(result.ApprovalDate),
	}

	return
}

// ApproveLoan handles the approval of a loan proposal.
// nolint
func (s *srv) InvestInLoan(ctx context.Context, req *gen.InvestInLoanRequest) (res *gen.InvestInLoanResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.InvestInLoanRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Extract claims from context.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate role_id from claims.
	if !slices.Contains(constant.AllowedRolesInvestLoan, claims.RoleId) {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Set InvestorId from claims.
	param.InvestorId = claims.UserId

	// Begin core process for the request.
	result, err := s.service.InvestInLoan(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.InvestInLoanResponse{
		LoanId:         result.LoanId,
		InvestedAmount: result.InvestedAmount,
		State:          result.State,
	}

	return
}

// nolint
func (s *srv) DisburseLoan(ctx context.Context, req *gen.DisburseLoanRequest) (res *gen.DisburseLoanResponse, err error) {
	// Cast and validate request.
	param := util.CastStruct[model.DisburseLoanRequest](req)
	if err = util.ValidateStruct(param); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Extract claims from context.
	claims, ok := util.ClaimsFromContext(ctx)
	if !ok {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Validate role_id from claims.
	if !slices.Contains(constant.AllowedRolesDisburseLoan, claims.RoleId) {
		err = constant.ErrPermissionDenied
		util.LogContext(ctx).Warn(err.Error())
		return
	}

	// Set OfficerId from claims.
	param.OfficerId = claims.UserId

	// Begin core process for the request.
	result, err := s.service.DisburseLoan(ctx, param)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &gen.DisburseLoanResponse{
		LoanId:           result.LoanId,
		DisbursementDate: timestamppb.New(result.DisbursementDate),
		State:            result.State,
	}

	return
}
