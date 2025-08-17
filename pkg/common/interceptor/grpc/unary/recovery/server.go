package recovery

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				debug.PrintStack()
				l.Error(fmt.Sprintf("%s PANIC: %v", util.GetMethod(info.FullMethod), r))
				err = status.Error(codes.Internal, "Internal error")
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false
		return resp, err
	}
}
