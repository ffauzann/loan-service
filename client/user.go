package client

import (
	"context"

	"github.com/ffauzann/loan-service/proto/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *userClient) CloseAccount(ctx context.Context, req *emptypb.Empty) (*gen.CloseAccountResponse, error) {
	return c.userClient.CloseAccount(ctx, req)
}
