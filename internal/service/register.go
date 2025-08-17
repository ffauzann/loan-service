package service

import (
	"context"
	"database/sql"
	"sync"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/jmoiron/sqlx"
)

func (s *service) Register(ctx context.Context, req *model.RegisterRequest) (res *model.RegisterResponse, err error) {
	// Begin tx.
	tx, err := s.repository.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}
	defer func() { s.repository.db.EndTx(ctx, tx, err) }()

	// Validate user existence.
	isUserExist, err := s.IsUserExist(ctx, &model.IsUserExistRequest{
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Return failed with reason if user already exists.
	if isUserExist.IsExist {
		res = &model.RegisterResponse{
			StatusCode: constant.RSCFailed,
			Reasons:    isUserExist.Reasons,
		}
		return
	}

	// Begin to register new user.
	_, err = s.createUser(ctx, tx, req)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.RegisterResponse{
		StatusCode: constant.RSCSucceed,
	}

	return
}

func (s *service) createUser(ctx context.Context, tx *sqlx.Tx, req *model.RegisterRequest) (user *model.User, err error) {
	// Prepare concurrent for hashing due it could take quite some times
	var wg sync.WaitGroup
	chErr := make(chan error, 2) //nolint
	fnHash := func(pwd string) {
		defer wg.Done()
		hashedPassword, err := util.HashPassword(pwd)
		if err != nil {
			chErr <- err
			return
		}
		req.User.UserPassword = hashedPassword
	}

	// Begin concurrent
	wg.Add(1)
	// go fnHash(s.config.Encryption.MasterPassword, constant.MasterPasswordType)
	go fnHash(req.User.PlainPassword)
	wg.Wait()

	// Begin non-blocking read channel
	select {
	case err = <-chErr: // Error occurred
		util.LogContext(ctx).Error(err.Error())
		return
	default: // No error, moving on
	}

	user = util.CastStruct[model.User](req.User)
	user.Password = []byte(req.User.UserPassword)
	user.Status = constant.UserStatusActive
	if err = s.repository.db.CreateUser(ctx, user, tx); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
