package mw

import (
	"context"
	"example/comments/internal/trace"

	"google.golang.org/grpc"
)

func Trace(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctxS, span := trace.Tracer().
		Start(
			ctx,
			info.FullMethod,
		)
	defer span.End()
	resp, err = handler(ctxS, req)
	return
}
