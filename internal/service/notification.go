package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
)

func (s *service) notifyLoanFullyFunded(ctx context.Context, loanId uint64) (err error) {
	investments, err := s.repository.db.GetInvestmentsByLoanId(ctx, loanId, nil)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}
	if len(investments) == 0 {
		util.LogContext(ctx).Info(fmt.Sprintf("No investments found for loan ID %d", loanId))
		return
	}

	// Get user details for each investment
	var userIds []uint64
	for _, investment := range investments {
		if slices.Contains(userIds, investment.InvestorId) {
			continue // Skip if investor already processed.
		} else {
			// Append unique investor ID to the list.
			userIds = append(userIds, investment.InvestorId)
		}
	}

	// Fetch user details for all unique investor IDs.
	var investors []*model.User
	investors, err = s.repository.db.GetUserByIds(ctx, userIds, nil)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Publish notification to each investor.
	for _, investor := range investors {
		if err = s.repository.messaging.Publish(ctx, &model.Message{
			Topic: constant.TopicFullyInvested,
			Payload: &model.EmailRequest{
				From:    "loan@service.com",
				To:      investor.Email,
				Subject: "Investment Confirmation",
				Body:    "You have successfully invested",
			},
		}); err != nil {
			util.LogContext(ctx).Error(err.Error())
			return
		}
	}

	return
}

func (s *service) SendMail(ctx context.Context, req *model.EmailRequest) (err error) {
	return s.repository.notification.SendMail(ctx, req)
}
