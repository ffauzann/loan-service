package service

import (
	"context"
	"time"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/repository"

	"go.uber.org/zap"
)

type Service interface {
	AuthService
	UserService
	LoanService
}

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (res *model.RegisterResponse, err error)
	IsUserExist(ctx context.Context, req *model.IsUserExistRequest) (res *model.IsUserExistResponse, err error)
	Login(ctx context.Context, req *model.LoginRequest) (res *model.LoginResponse, err error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (res *model.RefreshTokenResponse, err error)
	Jwks(ctx context.Context) (jwks []*model.Jwk, err error)
}

type UserService interface {
	CloseAccount(ctx context.Context, req *model.CloseAccountRequest) (res *model.CloseAccountResponse, err error)
}

type LoanService interface {
	CreateLoan(ctx context.Context, req *model.CreateLoanRequest) (res *model.CreateLoanResponse, err error)
	ApproveLoan(ctx context.Context, req *model.ApproveLoanRequest) (res *model.ApproveLoanResponse, err error)
	InvestInLoan(ctx context.Context, req *model.InvestInLoanRequest) (res *model.InvestInLoanResponse, err error)
	DisburseLoan(ctx context.Context, req *model.DisburseLoanRequest) (res *model.DisburseLoanResponse, err error)
}

type service struct {
	config     *model.AppConfig
	logger     *zap.Logger
	repository repositoryWrapper
}

type repositoryWrapper struct {
	db           repository.DBRepository
	redis        repository.RedisRepository
	notification repository.NotificationRepository
}

func New(db repository.DBRepository, redis repository.RedisRepository, notif repository.NotificationRepository, config *model.AppConfig, logger *zap.Logger) Service {
	return &service{
		config: config,
		logger: logger,
		repository: repositoryWrapper{
			db:           db,
			redis:        redis,
			notification: notif,
		},
	}
}

var now = time.Now // For mocking purpose later.
