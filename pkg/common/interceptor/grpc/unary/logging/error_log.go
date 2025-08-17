package logging

import (
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type errorLog struct {
	errorResponse
	IsKnown bool `json:"is_known"`
}

func logError(l *zap.Logger, msg string, ok bool, code codes.Code, err error) {
	eLog := errorLog{
		IsKnown: ok,
		errorResponse: errorResponse{
			Code:    code,
			Message: err.Error(),
		},
	}
	l.Info(msg, zap.Any("error", util.StructToMap(eLog)))
}
