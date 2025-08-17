package kafka

import (
	"context"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/service"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type srv struct {
	service service.Service
}

func New(consumer *kafka.Reader, svc service.Service) {
	s := &srv{
		service: svc,
	}

	for {
		ctx := context.Background()
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			util.Log().Error("failed to read kafka msg", zap.Error(err))
			continue
		}

		switch msg.Topic {
		case constant.TopicFullyInvested:
			go func() {
				if err := s.FullyInvested(ctx, msg.Value); err != nil {
					util.Log().Error(err.Error())
				}
			}()
		default:
			util.Log().Warn("unhandled kafka topic", zap.String("topic", msg.Topic))
		}
	}
}
