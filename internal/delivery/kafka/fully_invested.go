package kafka

import (
	"context"
	"encoding/json"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
)

func (s *srv) FullyInvested(ctx context.Context, msg []byte) (err error) {
	req := &model.EmailRequest{}
	if err = json.Unmarshal(msg, req); err != nil {
		util.Log().Error(err.Error())
		return
	}

	if err = s.service.SendMail(ctx, req); err != nil {
		util.Log().Error(err.Error())
		return
	}

	return
}
