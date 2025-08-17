package constant

import (
	"errors"

	"google.golang.org/grpc/codes"
)

// Known gRPC errors.
var (
	// Generic errors.
	ErrNotFound         = errors.New("Not found")
	ErrNoArg            = errors.New("No argument given")
	ErrInternal         = errors.New("Internal error")
	ErrPermissionDenied = errors.New("Permission denied")
	ErrUnauthenticated  = errors.New("Unauthenticated")

	// Specific errors.
	ErrInvalidMethod              = errors.New("Invalid method")
	ErrInvalidUsernamePassword    = errors.New("Invalid username/password")
	ErrPasswordIsTooWeak          = errors.New("Password is too weak")
	ErrMalformedEmail             = errors.New("Malformed email")
	ErrInvalidUserIdType          = errors.New("Invalid user ID type")
	ErrUserNotFound               = errors.New("User not found")
	ErrUserIsNotActive            = errors.New("User is blocked/closed")
	ErrUserAlreadyExists          = errors.New("User already exists")
	ErrInvalidToken               = errors.New("Invalid/expired token")
	ErrUnspecifiedAction          = errors.New("Unspecified Action")
	ErrLoanNotApproved            = errors.New("Loan not approved")
	ErrInvestmentAmountOutOfRange = errors.New("Investment amount out of range")
	ErrLoanNotProposed            = errors.New("Loan not proposed")
	ErrLoanNotFullyInvested       = errors.New("Loan not fully invested")
)

// All client-safe errors goes here.
var (
	MapGRPCErrCodes = map[error]codes.Code{
		// For HTTP mapping: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
		ErrInvalidMethod:              codes.InvalidArgument,
		ErrInvalidUsernamePassword:    codes.InvalidArgument,
		ErrMalformedEmail:             codes.InvalidArgument,
		ErrInvalidUserIdType:          codes.InvalidArgument,
		ErrUnspecifiedAction:          codes.FailedPrecondition,
		ErrPasswordIsTooWeak:          codes.FailedPrecondition,
		ErrNoArg:                      codes.FailedPrecondition,
		ErrLoanNotApproved:            codes.FailedPrecondition,
		ErrInvestmentAmountOutOfRange: codes.FailedPrecondition,
		ErrLoanNotProposed:            codes.FailedPrecondition,
		ErrLoanNotFullyInvested:       codes.FailedPrecondition,
		ErrNotFound:                   codes.NotFound,
		ErrUserNotFound:               codes.NotFound,
		ErrUserAlreadyExists:          codes.AlreadyExists,
		ErrPermissionDenied:           codes.PermissionDenied,
		ErrUserIsNotActive:            codes.PermissionDenied,
		ErrInternal:                   codes.Internal,
		ErrInvalidToken:               codes.Unauthenticated,
		ErrUnauthenticated:            codes.Unauthenticated,
	}
)
