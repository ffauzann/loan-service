package logging

import (
	"context"
	"time"

	"github.com/ffauzann/loan-service/pkg/common/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(l *zap.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		method := util.GetMethod(info.FullMethod)
		resp, err = handler(ctx, req)
		msg := formatMessage(start, method)

		if err != nil {
			err = func(err error) error {
				code, err, ok := o.getError(err)
				logError(o.prepareLog(ctx, l, resp), msg, ok, code, err)

				// Handle known error.
				if ok {
					return &errorResponse{
						Code:    code,
						Message: err.Error(),
					}
				}

				// Handle unknown error.
				return &errorResponse{
					Code:    code,
					Message: err.Error(),
				}
			}(err)

			return
		}

		// Only log success requests.
		o.prepareLog(ctx, l, resp).Info(msg)

		return
	}
}
