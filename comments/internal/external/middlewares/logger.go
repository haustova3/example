package mwc

import (
	"context"
	"example/comments/internal/logger"

	"google.golang.org/grpc"
)

func Logger(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, _ ...grpc.CallOption) error {
	logger.Infow(ctx, "SENT grpc request", "method", method)
	err := invoker(ctx, method, req, reply, cc)
	if err != nil {
		logger.Warnw(ctx, "GRPC request failed", "error", err)
	}
	return err
}
