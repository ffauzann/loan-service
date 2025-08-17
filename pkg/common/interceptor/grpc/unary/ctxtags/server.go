package ctxtags

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		reqMD, _ := metadata.FromIncomingContext(ctx)

		t := NewTags().
			Set(CIDKey, uuid.New().String()).
			Set(MDKey, reqMD).
			Set(ReqKey, req)

		ctx = SetInContext(ctx, t)

		return handler(ctx, req)
	}
}
