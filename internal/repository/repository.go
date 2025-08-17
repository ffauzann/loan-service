package repository

import (
	"context"
	"database/sql"
	"net/smtp"
	"time"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewDB(db *sqlx.DB, config *model.AppConfig, logger *zap.Logger) DBRepository {
	return &dbRepository{
		db: db,
		common: common{
			config: config,
			logger: logger,
		},
	}
}

func NewRedis(client *redis.Client, config *model.AppConfig, logger *zap.Logger) RedisRepository {
	return &redisRepository{
		redis: client,
		common: common{
			config: config,
			logger: logger,
		},
	}
}

func NewNotification(client *smtp.Client, config *model.AppConfig, logger *zap.Logger) NotificationRepository {
	return &notificationRepository{
		smtp: client,
		common: common{
			config: config,
			logger: logger,
		},
	}
}

type DBRepository interface {
	DBTxRepository
	DBUserRepository
	DBLoanRepository
}

type DBTxRepository interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (tx *sqlx.Tx, err error)
	EndTx(ctx context.Context, tx *sqlx.Tx, err error)
}

type DBUserRepository interface {
	CreateUser(ctx context.Context, user *model.User, tx *sqlx.Tx) error
	IsUserExist(ctx context.Context, userIdType constant.UserIdType, userIdVal string) (isExist bool, err error)
	GetUserByOneOfIdentifier(ctx context.Context, val string) (user *model.User, err error)
	CloseAccount(ctx context.Context, req *model.CloseAccountRequest, tx *sqlx.Tx) (err error)

	GetUserByIds(ctx context.Context, userIds []uint64, tx *sqlx.Tx) (users []*model.User, err error)
}

type DBLoanRepository interface {
	CreateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) error
	GetLoanById(ctx context.Context, loanID uint64, tx *sqlx.Tx) (loan *model.Loan, err error)
	UpdateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) (err error)
	ApproveLoan(ctx context.Context, approval *model.LoanApproval, tx *sqlx.Tx) error
	CreateLoanInvestment(ctx context.Context, investment *model.LoanInvestment, tx *sqlx.Tx) error
	CreateLoanDisbursement(ctx context.Context, disbursement *model.LoanDisbursement, tx *sqlx.Tx) (err error)

	GetInvestmentsByLoanId(ctx context.Context, loanId uint64, tx *sqlx.Tx) ([]*model.LoanInvestment, error)
}

type RedisRepository interface {
	RegisterUserDevice(ctx context.Context, deviceId string, token *model.Token) error
}

type NotificationRepository interface {
	SendMail(ctx context.Context, req *model.EmailRequest) error
}

type common struct {
	config *model.AppConfig
	logger *zap.Logger
}

type dbRepository struct {
	db *sqlx.DB
	common
}

type redisRepository struct {
	redis *redis.Client
	common
}

type notificationRepository struct {
	smtp *smtp.Client
	common
}

var now = time.Now // For mocking purpose later.
