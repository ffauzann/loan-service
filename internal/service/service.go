package service

import (
	"context"
	"time"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/repository"

	"go.uber.org/zap"
)

type Service interface {
	AuthService
	UserService
	LoanService
	NotificationService
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

type NotificationService interface {
	SendMail(ctx context.Context, req *model.EmailRequest) (err error)
}

type service struct {
	config     *model.AppConfig
	logger     *zap.Logger
	repository repositoryWrapper
}

type repositoryWrapper struct {
	db           repository.DBRepository
	redis        repository.RedisRepository
	messaging    repository.MessagingRepository
	notification repository.NotificationRepository
}

func New(db repository.DBRepository, redis repository.RedisRepository, messaging repository.MessagingRepository, notif repository.NotificationRepository, config *model.AppConfig, logger *zap.Logger) Service {
	return &service{
		config: config,
		logger: logger,
		repository: repositoryWrapper{
			db:           db,
			redis:        redis,
			messaging:    messaging,
			notification: notif,
		},
	}
}

var now = time.Now // For mocking purpose later.
