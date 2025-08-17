package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	commonUtil "github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/constant"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
)

func (s *service) IsUserExist(ctx context.Context, req *model.IsUserExistRequest) (res *model.IsUserExistResponse, err error) {
	// Prepare concurrent
	var wg sync.WaitGroup
	chErr := make(chan error, 2)     //nolint
	chReason := make(chan string, 2) //nolint
	fnIsExist := func(userIdType constant.UserIdType, userIdVal string) {
		defer wg.Done()
		isExist, err := s.repository.db.IsUserExist(ctx, userIdType, userIdVal)
		if err != nil {
			chErr <- err
			return
		}

		if isExist {
			chReason <- fmt.Sprintf("user with %s %s already exist", strings.Replace(string(userIdType), "_", " ", 1), userIdVal)
		}
	}

	// Begin concurrent
	if req.Email != "" {
		wg.Add(1)
		go fnIsExist(constant.UserIdTypeEmail, req.Email)
	}
	if req.PhoneNumber != "" {
		wg.Add(1)
		go fnIsExist(constant.UserIdTypePhoneNumber, req.PhoneNumber)
	}
	wg.Wait()

	// Begin non-blocking read channel
	if err = commonUtil.ErrorFromChannel(chErr); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Read all channel values if any
	close(chReason)
	reasons := []string{}
	for s := range chReason {
		reasons = append(reasons, s)
	}

	// Format response
	res = &model.IsUserExistResponse{
		IsExist: len(reasons) > 0,
		Reasons: reasons,
	}

	return
}

func (s *service) CloseAccount(ctx context.Context, req *model.CloseAccountRequest) (res *model.CloseAccountResponse, err error) {
	// Update user status to CLOSED regardless of its original status.
	err = s.repository.db.CloseAccount(ctx, req, nil)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.CloseAccountResponse{
		UserId: req.UserId,
		Status: constant.UserStatusClosed,
	}

	return
}
