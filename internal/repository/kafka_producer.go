package repository

import (
	"context"
	"encoding/json"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
	"github.com/segmentio/kafka-go"
)

func (r *messagingRepository) Publish(ctx context.Context, msg *model.Message) (err error) {
	bPayload, err := json.Marshal(msg.Payload)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	err = r.producer.WriteMessages(ctx, kafka.Message{
		Topic: msg.Topic,
		Key:   []byte("msg"),
		Value: bPayload,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	return
}
